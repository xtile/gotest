package main

//TODO:
// run 3 websocket connections with goroutines
// store current price in variables
// every n sec compare prices log prices and their diff

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"

	"github.com/sacOO7/gowebsocket"
)

func startWebSocketDataTransfer(exchange string) {
	//exchange := "HUOBI"
	//exchange = "OKEX"
	//exchange = "BINANCE"

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	socket := gowebsocket.New("")
	if exchange == "BINANCE" {
		// binance
		socket = gowebsocket.New("wss://stream.binance.com:9443/ws/BTCUSDT@kline_1m")
	} else if exchange == "HUOBI" {
		// huobi
		socket = gowebsocket.New("wss://api.huobi.pro/ws")

	} else if exchange == "OKEX" {
		// okex
		socket = gowebsocket.New("wss://real.okex.com:8443/ws/v3")
	}

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println(exchange + ": Connected to server")

		if exchange == "BINANCE" {
			//binance
			//socket.SendText("{  \"method\": \"SUBSCRIBE\",  \"params\": [    \"btcusdt@kline_1m\" ],  \"id\": 1}")
			socket.SendText("{  \"method\": \"SUBSCRIBE\",  \"params\": [    \"btcusdt@trade\" ],  \"id\": 1}")
		} else if exchange == "HUOBI" {
			//huobi
			socket.SendText("{ \"sub\": \"market.btcusdc.ticker\", \"id\": \"1\" }")
		} else if exchange == "OKEX" {
			//okex
			socket.SendText("{\"op\": \"subscribe\", \"args\": [\"spot/ticker:BTC-USDT\"]}")
		}

	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		//log.Println(exchange + ": Recieved message " + message)

		if exchange == "BINANCE" {
			// binance - no action, because we get text data from binance

			var result map[string]string
			json.Unmarshal([]byte(message), &result)

			p := result["p"]
			log.Println(exchange + ": price " + p)

			return

		}
		log.Println(exchange + ": Recieved message " + message)

	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		//log.Println(exchange+": Recieved binary data ", data)

		if exchange == "BINANCE" {
			// binance - no action, because we get text data from binance

		} else if exchange == "HUOBI" {
			// huobi
			//strDataOriginal := string(data)
			//log.Println("original data ", strDataOriginal)

			bytesZipped := bytes.NewReader(data[:])
			zipReader, err := gzip.NewReader(bytesZipped)
			if err != nil {
				log.Println("error1 ", err)
				return
			}
			defer zipReader.Close()
			bytesUnzipped, err := ioutil.ReadAll(zipReader)
			if err != nil {
				log.Println("error2 ", err)
				return
			}
			strUnzipped := string(bytesUnzipped)

			var result map[string]interface{}
			json.Unmarshal([]byte(strUnzipped), &result)

			if result["tick"] != nil {
				p := result["tick"].(map[string]interface{})
				if p != nil {
					price := p["lastPrice"].(float64)
					//strPrice, _ :=  strconv.ParseFloat(price, 64)
					strPrice := strconv.FormatFloat(price, 'f', 2, 64)
					log.Println(exchange + ": price " + strPrice)
				}
			} else {
				log.Println(exchange+": decoded message:  ", strUnzipped)

			}
			//p := result["tick"]
			//log.Println(exchange + ": price " + p)

		} else if exchange == "OKEX" {
			// okex
			// deflating compressed binary data from OKex
			var b bytes.Buffer
			r := flate.NewReader(bytes.NewReader(data))
			b.ReadFrom(r)
			r.Close()
			strData := string(b.Bytes())
			//log.Println(exchange+": inflated data ", strData)

			var result map[string]interface{} //json.RawMessage
			json.Unmarshal([]byte(strData), &result)

			if result["data"] != nil {
				p := result["data"].([]interface{})
				if p != nil && p[0] != nil {
					pp := p[0].(map[string]interface{})
					if pp != nil {
						price := pp["last"].(string)
						log.Println(exchange + ": price " + price)
					}
				}
			}
		}

	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println(exchange + ": Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println(exchange + ": Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println(exchange + ": Disconnected from server ")
		return
	}

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}

func main() {

	runtime.GOMAXPROCS(3)

	var wg sync.WaitGroup
	wg.Add(3)

	//go startWebSocketDataTransfer("BINANCE")
	go startWebSocketDataTransfer("HUOBI")
	go startWebSocketDataTransfer("OKEX")

	log.Println("Waiting To Finish")
	wg.Wait()

	log.Println("\nTerminating Program")

}

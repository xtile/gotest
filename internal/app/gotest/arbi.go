package gotest

import (
	"github.com/sirupsen/logrus"

	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/sacOO7/gowebsocket"
)

type ArbiLogger struct {
	config *Config
	logger *logrus.Logger
}

func New(s *Config) *ArbiLogger {
	return &ArbiLogger{
		config: s,
		logger: logrus.New(),
	}
}

func (s *ArbiLogger) Start() error {

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("starting ARBI logger")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	runtime.GOMAXPROCS(4)

	var wg sync.WaitGroup
	wg.Add(4)

	go startWebSocketDataTransfer("BINANCE")
	go startWebSocketDataTransfer("HUOBI")
	go startWebSocketDataTransfer("OKEX")
	go comparePrices()
	go func() {
		for {
			sig := <-sigs
			fmt.Println(" WAITING SIG ")
			fmt.Println(sig)
			fmt.Println("WAITING SIG - AFTER SIG")
			done <- true
			fmt.Println("CLEANUP")
			os.Exit(1)
			fmt.Println("AFTER EXIT")

		}
	}()

	//c := make(chan os.Signal)
	//signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//go func() {
	//	<-c
	//	cleanup()
	//	os.Exit(1)
	//}()

	log.Println("awaiting signal")
	<-done
	log.Println("exiting")

	log.Println("Waiting To Finish")
	wg.Wait()

	log.Println("\nTerminating Program")

	return nil

}

func (s *ArbiLogger) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)

	if err != nil {
		return err
	}
	s.logger.SetLevel(level)

	return nil
}

var priceBinance, priceOKex, priceHuobi float64 = 0, 0, 0
var tsBinance, tsOKex, tsHuobi int = 0, 0, 0
var delta float64 = 0.002
var index int = 0
var avgPrice float64 = 0

//var BINANCE_MARKER string = "ethusdt@trade"
//var HUOBI_MARKER string = "ethusdc"
//var OKEX_MARKER string = "ETH-USDT"

var BINANCE_MARKER string = "ltcusdt@trade"
var HUOBI_MARKER string = "ltcusdt"
var OKEX_MARKER string = "LTC-USDT"

func comparePrices() {
	for {

		if priceBinance == 0 || priceHuobi == 0 || priceOKex == 0 {
			continue
		}

		index += 1

		avgPrice = (priceBinance + priceOKex + priceHuobi) / 3

		if index%10 == 2 {
			log.Println("Binance: ", priceBinance, " HUOBI: ", priceHuobi, "OKEX: ", priceOKex)
		}

		time.Sleep(2 * time.Second)
		if priceBinance-priceOKex > delta*avgPrice {
			log.Println("BINANCE > OKEX --------------------------------------------")
		}
		if priceBinance-priceHuobi > delta*avgPrice {
			log.Println("BINANCE > HUOBI ", priceBinance-priceHuobi)
		}
		if priceOKex-priceBinance > delta*avgPrice {
			log.Println("OKEX > BINANCE --------------------------------------------")
		}
		if priceOKex-priceHuobi > delta*avgPrice {
			log.Println("OKEX > HUOBI ", priceOKex-priceHuobi)
		}
		if priceHuobi-priceOKex > delta*avgPrice {
			log.Println("HUOBI > OKEX -", priceHuobi-priceOKex)
		}
		if priceHuobi-priceBinance > delta*avgPrice {
			log.Println("HUOBI > BINANCE:", priceHuobi-priceBinance)
		}
	}
}

func startWebSocketDataTransfer(exchange string) {

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
			//socket.SendText("{  \"method\": \"SUBSCRIBE\",  \"params\": [    \"btcusdt@trade\" ],  \"id\": 1}")
			socket.SendText("{  \"method\": \"SUBSCRIBE\",  \"params\": [    \"" + BINANCE_MARKER + "\" ],  \"id\": 1}")
		} else if exchange == "HUOBI" {
			//huobi
			//socket.SendText("{ \"sub\": \"market.btcusdc.ticker\", \"id\": \"1\" }")
			socket.SendText("{ \"sub\": \"market." + HUOBI_MARKER + ".ticker\", \"id\": \"1\" }")
		} else if exchange == "OKEX" {
			//okex
			//socket.SendText("{\"op\": \"subscribe\", \"args\": [\"spot/ticker:BTC-USDT\"]}")
			socket.SendText("{\"op\": \"subscribe\", \"args\": [\"spot/ticker:" + OKEX_MARKER + "\"]}")
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
			//PRICE log.Println(exchange + ": price " + p)
			priceBinance, _ = strconv.ParseFloat(p, 2)

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
					strPrice := strconv.FormatFloat(price, 'f', 2, 64)
					//PRICE log.Println(exchange + ": price " + strPrice)
					priceHuobi, _ = strconv.ParseFloat(strPrice, 2)

				}
			} else if result["ping"] != nil {
				pingData := result["ping"].(float64)
				pongData := "{ \"pong\": " + fmt.Sprintf("%.0f", pingData) + " }"
				//log.Println(exchange+": pongData:  ", pongData)

				socket.SendText(pongData)

			} else {
				log.Println(exchange+": decoded message:  ", strUnzipped)

			}

		} else if exchange == "OKEX" {
			// okex
			// deflating compressed binary data from OKex
			var b bytes.Buffer
			r := flate.NewReader(bytes.NewReader(data))
			b.ReadFrom(r)
			r.Close()
			strData := string(b.Bytes())

			var result map[string]interface{} //json.RawMessage
			json.Unmarshal([]byte(strData), &result)

			if result["data"] != nil {
				p := result["data"].([]interface{})
				if p != nil && p[0] != nil {
					pp := p[0].(map[string]interface{})
					if pp != nil {
						price := pp["last"].(string)
						//PRICE log.Println(exchange + ": price " + price)
						priceOKex, _ = strconv.ParseFloat(price, 2)
					}
				}
			}
		}

	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		//log.Println(exchange + ": Recieved ping " + data)
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
			log.Println("exchange: interrupt")
			socket.Close()
			return
		}
	}
}

package main

//TODO: 
// run 3 websocket connections with goroutines
// store current price in variables 
// every n sec compare prices log prices and their diff


import (
	"bytes"
	"compress/flate"
	"log"
	"os"
	"os/signal"

	"github.com/sacOO7/gowebsocket"
)

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// binance
	//socket := gowebsocket.New("wss://stream.binance.com:9443/ws/BTCUSDT@kline_1m")
	// okex
	socket := gowebsocket.New("wss://real.okex.com:8443/ws/v3")

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
		//binance
		//socket.SendText("{  \"method\": \"SUBSCRIBE\",  \"params\": [    \"btcusdt@kline_1m\" ],  \"id\": 1}")
		//socket.SendText("{  \"method\": \"SUBSCRIBE\",  \"params\": [    \"btcusdt@trade\" ],  \"id\": 1}")
		//okex
		socket.SendText("{\"op\": \"subscribe\", \"args\": [\"spot/ticker:BTC-USDT\"]}")

		/*
			{  "method": "SUBSCRIBE",  "params": [    "btcusdt@aggTrade",    "btcusdt@depth"  ],  "id": 1}
		*/
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		log.Println("Recieved binary data ", data)
		// deflating compressed binary data from OKex
		var b bytes.Buffer
		r := flate.NewReader(bytes.NewReader(data))
		b.ReadFrom(r)
		r.Close()
		strData := string(b.Bytes())
		log.Println("inflated data ", strData)

	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
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

package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sacOO7/gowebsocket"
)

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New("wss://stream.binance.com:9443/ws/BTCUSDT@kline_1m")

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
		//socket.SendText("{  \"method\": \"SUBSCRIBE\",  \"params\": [    \"btcusdt@kline_1m\" ],  \"id\": 1}")
		socket.SendText("{  \"method\": \"SUBSCRIBE\",  \"params\": [    \"btcusdt@trade\" ],  \"id\": 1}")

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

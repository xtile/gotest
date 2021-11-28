module github.com/xtile/gotest

go 1.17

//replace github.com/xtile/gotest => ./
replace github.com/xtile/gotest => ./ //github.com/xtile/arbi@latest

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/sacOO7/gowebsocket v0.0.0-20210515122958-9396f1a71e23
	github.com/sirupsen/logrus v1.8.1
//github.com/xtile/gotest@latest
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/sacOO7/go-logger v0.0.0-20180719173527-9ac9add5a50d // indirect
	github.com/xtile/gotest v0.0.0-20211127234445-1842d1f2eb79 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)

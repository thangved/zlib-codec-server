package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	encryption "thangved.com/zlib-codec-server"

	"go.temporal.io/sdk/converter"
)

func main() {
	EnvPort, _ := strconv.Atoi(os.Getenv("PORT"))

	EnvKeyID := os.Getenv("KEY_ID")

	flag.Parse()

	handler := converter.NewPayloadCodecHTTPHandler(&encryption.Codec{
		KeyID: EnvKeyID,
	}, converter.NewZlibCodec(converter.ZlibCodecOptions{AlwaysEncode: true}))

	srv := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(EnvPort),
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	log.Println("Server started at", srv.Addr+":"+strconv.Itoa(EnvPort))

	select {
	case <-sigCh:
		_ = srv.Close()
	case err := <-errCh:
		log.Fatal(err)
	}

	log.Println("codec-server stopped")
}

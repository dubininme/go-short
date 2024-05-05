package main

import (
	"flag"
	"log"
	"net/http"
	"net/rpc"

	"github.com/dubininme/go-short/key"
	"github.com/dubininme/go-short/shortener"
	"github.com/dubininme/go-short/storage"
)

var (
	keyChar    = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	listenAddr = flag.String("http", ":3000", "http listen addres")
	hostname   = flag.String("http", "localhost:3000", "http host")
	dataFile   = flag.String("file", "store.json", "data store filename")
	mainAddr   = flag.String("master", "", "RPC main address")
	rpcEnabled = flag.Bool("rpc", false, "enable RPC server")
)

func main() {
	flag.Parse()

	keyGen := key.NewKeyGenerator(keyChar)
	var st shortener.Storage
	st = storage.NewUrlStorage(keyGen, *dataFile)
	if *mainAddr != "" {
		st = storage.NewProxyStorage(st.(*storage.UrlStorage), *mainAddr)
	}

	shortener := shortener.NewUrlShortnener(st)
	if *rpcEnabled {
		rpc.RegisterName("Shortener", shortener)
		rpc.HandleHTTP()
	}

	http.HandleFunc("/", shortener.Redirect)
	http.HandleFunc("/add", shortener.Add)
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatal("Failed to start server %v", err)
	}
}

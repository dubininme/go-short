package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dubininme/go-short/key"
	"github.com/dubininme/go-short/shortener"
	"github.com/dubininme/go-short/storage"
)

var keyChar = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var dataFile = flag.String("file", "store.json", "data store filename")

func main() {

	if dataFile == nil {
		log.Fatal("Data file not initialized")
	}

	keyGen := key.NewKeyGenerator(keyChar)
	storage := storage.NewUrlStorage(keyGen, *dataFile)
	shortener := shortener.NewUrlShortnener(storage)

	http.HandleFunc("/", shortener.Redirect)
	http.HandleFunc("/add", shortener.Add)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Failed to start server %v", err)
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/dubininme/go-short/key"
	"github.com/dubininme/go-short/shortener"
	"github.com/dubininme/go-short/storage"
)

var keyChar = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {
	keyGen := key.NewKeyGenerator(keyChar)
	storage := storage.NewUrlStorage(keyGen, "storage.gob")
	shortener := shortener.NewUrlStorage(storage)

	http.HandleFunc("/", shortener.Redirect)
	http.HandleFunc("/add", shortener.Add)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Failed to start server %v", err)
	}
}

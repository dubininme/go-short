package storage

import (
	"log"
	"net/rpc"
)

type ProxyStorage struct {
	storage *UrlStorage
	client  *rpc.Client
}

func NewProxyStorage(storage *UrlStorage, addr string) *ProxyStorage {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatal("Error creating new proxy storage")
	}

	return &ProxyStorage{storage: storage, client: client}
}

func (s *ProxyStorage) Get(key, url *string) error {
	if err := s.storage.Get(key, url); err == nil {
		return nil
	}

	// rpc call to master:
	if err := s.client.Call("Store.Get", key, url); err != nil {
		return err
	}

	s.storage.Set(key, url) // update local cache
	return nil
}

func (s *ProxyStorage) Put(url, key *string) error {
	// rpc call to main
	if err := s.client.Call("Store.Put", url, key); err != nil {
		return err
	}

	// update local storage
	s.storage.Set(key, url)
	return nil
}

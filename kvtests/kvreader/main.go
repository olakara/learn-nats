package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, _ := nc.JetStream()

	kv, err := js.KeyValue("MY_BUCKET")

	if err != nil {
		log.Fatal(err)
	}

	entry, err := kv.Get("4a16c5bf-ee9f-4364-80f9-4a52b1004b79")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Key:", entry.Key())
	fmt.Println("Value:", string(entry.Value()))
	fmt.Println("Revision:", entry.Revision())
	fmt.Println("Created:", entry.Created())

	kv.Put("4a16c5bf-ee9f-4364-80f9-4a52b1004b79", []byte("simple value"))

	another, err := kv.Get("f86d3ad7-d0ad-4b0e-9f58-7ce9c2ebf4be")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Key:", another.Key())
	fmt.Println("Value:", string(another.Value()))
	fmt.Println("Revision:", another.Revision())
	fmt.Println("Created:", another.Created())
}

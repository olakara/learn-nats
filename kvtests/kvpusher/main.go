package main

import (
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

	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "MY_BUCKET",
	})

	if err != nil {
		log.Fatal(err)
	}

	kv.Put("4a16c5bf-ee9f-4364-80f9-4a52b1004b79", []byte("hello"))
	kv.PutString("f86d3ad7-d0ad-4b0e-9f58-7ce9c2ebf4be", "abdel!")
}

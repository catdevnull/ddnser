package main

import (
	"log"

	"nulo.in/ddnser/nameservers"
)

func main() {
	njalla := nameservers.Njalla{Key: "yourkey"}
	record, err := njalla.SetRecord("estoesprueba.nulo.in", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(record)
	record, err = njalla.SetRecord("estoesprueba.nulo.in", "1.1.1.1")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(record)

}

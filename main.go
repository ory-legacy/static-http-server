package main

import (
	"net/http"
	"encoding/json"
	"os"
	"log"
	"strings"
)

type Configuration struct {
	Port    string
}

func main() {
	file, err := os.Open("config.json")

	if err != nil {
		log.Fatal("Could not open config: ", err)
		return
	}

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)

	if err != nil {
		log.Fatal("Could not decode config: ", err)
		return
	}

	var port = []string{":", configuration.Port};
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(strings.Join(port, ""), nil))
}

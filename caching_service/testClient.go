package main

import (
	"net/http"
	"log"
	"fmt"
	"time"
)

func foo(val int) {
	//do something with val that causes global cache update
}

func main() {
	val, err := http.Get("http://localhost:8085/primary?val=2")
	if err != nil {
		log.Fatal(err)
		fmt.Errorf("Error requesting value from primary cache")
	}

	foo(val)

	//set a timeout for 5 seconds
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	//After calling function, retrieve another value from primary cache
	val, err = client.Get("http://localhost:8085/primary?val=3")
	if err != nil {
		//if error, try secondary cache
		val, err = http.Get("http://localhost:8085/secondary?val=3")
		if err != nil {
			log.Fatal(err)
			fmt.Errorf("Error requesting value from secondary cache")
		}

	}
}

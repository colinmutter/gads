package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	gads "github.com/Getsidecar/gads/v201710"
)

var configJson = flag.String("oauth", "./oauth.json", "API credentials")

func main() {
	flag.Parse()
	config, err := gads.NewCredentialsFromFile(*configJson)
	if err != nil {
		log.Fatal(err)
	}

	cs := gads.NewCustomerService(&config.Auth)

	customers, err := cs.GetCustomers()

	fmt.Println(customers)
	customersJSON, err := json.MarshalIndent(customers, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", customersJSON)

}

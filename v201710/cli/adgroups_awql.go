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
	// show all Campaigns
	cs := gads.NewAdGroupService(&config.Auth)
	fmt.Printf("AdGroups\n")

	foundAdGroups, totalCount, err := cs.Query("SELECT AdGroupId, Name, Status WHERE CampaignId = '123'")
	fmt.Println(totalCount)
	if err != nil {
		log.Fatal(err)
	}
	for _, adGroup := range foundAdGroups {
		adGroupJSON, _ := json.MarshalIndent(adGroup, "", "  ")
		fmt.Printf("%s\n", adGroupJSON)
	}

}

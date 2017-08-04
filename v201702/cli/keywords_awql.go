package main

import (
	"encoding/json"
	"flag"
	"fmt"
	gads "github.com/colinmutter/gads/v201609"
	"log"
)

var configJson = flag.String("oauth", "./oauth.json", "API credentials")

func main() {
	flag.Parse()
	config, err := gads.NewCredentialsFromFile(*configJson)
	if err != nil {
		log.Fatal(err)
	}

	// show all Keywords
	cs := gads.NewAdGroupCriterionService(&config.Auth)
	fmt.Printf("Keywords\n")
	foundKeywords, totalCount, err := cs.Query("SELECT Id, KeywordText, KeywordMatchType WHERE AdGroupId = '123'")
	fmt.Println(totalCount)
	if err != nil {
		log.Fatal(err)
	}
	for _, keyword := range foundKeywords {
		keywordJSON, _ := json.MarshalIndent(keyword, "", "  ")
		fmt.Printf("%s\n", keywordJSON)
	}

}

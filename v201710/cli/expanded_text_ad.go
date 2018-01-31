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

	aga := gads.NewAdGroupAdService(&config.Auth)

	ads, err := aga.Mutate(
		gads.AdGroupAdOperations{
			"ADD": {
				gads.ExpandedTextAd{
					AdGroupId:     1234567890,
					FinalUrls:     []string{"https://classdo.com/en"},
					Path1:         "path1",
					Path2:         "path2",
					HeadlinePart1: "test headline",
					HeadlinePart2: "test headline2",
					Description:   "test line one",
				},
			},
		},
	)

	fmt.Println(ads)
	adsJSON, err := json.MarshalIndent(ads, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", adsJSON)

}

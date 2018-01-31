package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	gads "github.com/Getsidecar/gads/v201710"
)

var configJson = flag.String("oauth", "./oauth.json", "API credentials")

// Your IDs should go here
var campaignId string = "1234567890"
var adGroupId string = "1234567890"
var criterionId string = "1234567890"

func main() {
	flag.Parse()
	config, err := gads.NewCredentialsFromFile(*configJson)
	if err != nil {
		log.Fatal(err)
	}

	service := gads.NewDataService(&config.Auth)

	// Test GetAdGroupBidLandscape

	adGroupBidLandscape, _, err := service.GetAdGroupBidLandscape(
		gads.Selector{
			Fields: []string{
				"AdGroupId",
				"Bid",
				"CampaignId",
				"LocalClicks",
				"LocalCost",
				"LocalImpressions",
				"PromotedImpressions",
				"StartDate",
				"EndDate",
			},
			Predicates: []gads.Predicate{
				{"AdGroupId", "EQUALS", []string{adGroupId}},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	for _, bidLandscape := range adGroupBidLandscape {
		landscapeJSON, _ := json.MarshalIndent(bidLandscape, "", "  ")
		fmt.Printf("%s\n", landscapeJSON)
	}

	// END Test GetAdGroupBidLandscape

	// Test QueryAdGroupBidLandscape

	foundBidLandscape, _, err := service.QueryAdGroupBidLandscape(fmt.Sprintf("SELECT AdGroupId, Bid, CampaignId, LocalClicks, LocalCost, LocalImpressions WHERE AdGroupId = '%v'", adGroupId))

	if err != nil {
		log.Fatal(err)
	}

	for _, bidLandscape := range foundBidLandscape {
		landscapeJSON, _ := json.MarshalIndent(bidLandscape, "", "  ")
		fmt.Printf("%s\n", landscapeJSON)
	}

	// END Test QueryAdGroupBidLandscape

	// Test GetCriterionBidLandscape

	criterionBidLandscape, _, err := service.GetCriterionBidLandscape(
		gads.Selector{
			Fields: []string{
				"AdGroupId",
				"Bid",
				"CampaignId",
				"CriterionId",
				"LocalClicks",
				"LocalCost",
				"LocalImpressions",
				"PromotedImpressions",
				"StartDate",
				"EndDate",
			},
			Predicates: []gads.Predicate{
				{"CriterionId", "EQUALS", []string{criterionId}},
				{"AdGroupId", "EQUALS", []string{adGroupId}},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	for _, bidLandscape := range criterionBidLandscape {
		landscapeJSON, _ := json.MarshalIndent(bidLandscape, "", "  ")
		fmt.Printf("%s\n", landscapeJSON)
	}

	// END Test GetCriterionBidLandscape

	// Test QueryCriterionBidLandscape

	queryCriterionBidLandscape, _, err := service.QueryCriterionBidLandscape(fmt.Sprintf("SELECT AdGroupId, Bid, CampaignId, CriterionId, LocalClicks, LocalCost, LocalImpressions WHERE AdGroupId = '%v' AND CriterionId = '%v'", adGroupId, criterionId))

	if err != nil {
		log.Fatal(err)
	}

	for _, bidLandscape := range queryCriterionBidLandscape {
		landscapeJSON, _ := json.MarshalIndent(bidLandscape, "", "  ")
		fmt.Printf("%s\n", landscapeJSON)
	}

	// END Test QueryCriterionBidLandscape

}

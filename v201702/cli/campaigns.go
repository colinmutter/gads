package main

import (
	"encoding/json"
	"fmt"
	"github.com/colinmutter/gads/v201609"
	"golang.org/x/oauth2"
	"log"
)

func main() {
	config, err := gads.NewCredentials(oauth2.NoContext)
	if err != nil {
		log.Fatal(err)
	}

	var pageSize int64 = 500
	var offset int64

	// show all Campaigns
	cs := gads.NewCampaignService(&config.Auth)
	paging := gads.Paging{
		Offset: offset,
		Limit:  pageSize,
	}
	fmt.Printf("\nCampaigns\n")
	for {
		foundCampaigns, totalCount, err := cs.Get(
			gads.Selector{
				Fields: []string{
					"Id",
					"Name",
					"Status",
					"ServingStatus",
					"StartDate",
					"EndDate",
					"AdServingOptimizationStatus",
					"Settings",
					"AdvertisingChannelType",
					"AdvertisingChannelSubType",
					"Labels",
					"TrackingUrlTemplate",
					"UrlCustomParameters",
				},
				Predicates: []gads.Predicate{
					{"Status", "EQUALS", []string{"PAUSED"}},
				},
				Ordering: []gads.OrderBy{
					{"Id", "ASCENDING"},
				},
				Paging: &paging,
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		for _, campaign := range foundCampaigns {
			campaignJson, _ := json.MarshalIndent(campaign, "", "  ")
			fmt.Printf("%s\n", campaignJson)
		}
		offset += pageSize
		paging.Offset = offset
		if totalCount < offset {
			break
		}
	}

}

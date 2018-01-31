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

	cs := gads.NewCampaignCriterionService(&config.Auth)

	//bidmodifier := float64(2)
	//var bidmodifier float64
	//bidmodifier = 1.04

	campaignCriterions := gads.CampaignCriterions{
		gads.NegativeCampaignCriterion{
			CampaignId: 123456789,
			Criterion: gads.Location{
				Id: 21168,
			},
		},
		gads.CampaignCriterion{
			CampaignId: 123456789,
			Criterion: gads.Location{
				Id: 21167,
			},
		},
	}

	criterions, err := cs.Mutate(
		gads.CampaignCriterionOperations{
			"ADD": campaignCriterions,
		},
	)

	fmt.Println(criterions)
	criterionJSON, _ := json.MarshalIndent(criterions, "", "  ")
	fmt.Printf("%s\n", criterionJSON)

	// show all Location Criterion
	fmt.Printf("Campaign Criterion\n")
	foundCriterions, totalCount, err := cs.Query("SELECT CampaignId,IsNegative,BidModifier,CriteriaType,Id,LocationName,DisplayType,ParentLocations WHERE CampaignId = '211793582' AND CriteriaType IN ['LOCATION']")
	fmt.Println(totalCount)
	if err != nil {
		log.Fatal(err)
	}
	for _, criterion := range foundCriterions {
		criterionJSON, _ := json.MarshalIndent(criterion, "", "  ")
		fmt.Printf("%s\n", criterionJSON)
	}

}

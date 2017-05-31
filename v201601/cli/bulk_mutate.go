package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	gads "github.com/bbachtel/gads/v201601"
)

var configJson = flag.String("oauth", "./oauth.json", "API credentials")

// Your campaign ID should go here
var campaignId int64 = 1234567890

func main() {
	flag.Parse()
	config, err := gads.NewCredentialsFromFile(*configJson)
	if err != nil {
		log.Fatal(err)
	}
	// Bulk Mutate
	ms := gads.NewMutateJobService(&config.Auth)
	policy := new(gads.BulkMutateJobPolicy)

	// If you need to add prerequisites
	//policy.PrerequisiteJobIds = append(policy.PrerequisiteJobIds, 123456)

	// Creating AdGroups
	/*ago := gads.AdGroupOperations{
		"ADD": {
			gads.AdGroup{
				Name:       "test ad group " + rand_str(10),
				Status:     "PAUSED",
				CampaignId: campaignId,
			},
			gads.AdGroup{
				Name:       "test ad group " + rand_str(10),
				Status:     "PAUSED",
				CampaignId: campaignId,
			},
		},
	}*/

	// Updating AdGroups
	/*ago := gads.AdGroupOperations{
		"SET": {
			gads.AdGroup{
				Id: 1234567890,
				CampaignId: campaignId,
				BiddingStrategyConfiguration: []gads.BiddingStrategyConfiguration{
					gads.BiddingStrategyConfiguration{
						StrategyType: "MANUAL_CPC",
						Bids: []gads.Bid{
							gads.Bid{
								Type:   "CpcBid",
								Amount: 2000000,
							},
						},
					},
				},
			},
			gads.AdGroup{
				Id: 1234567890,
				CampaignId: campaignId,
				BiddingStrategyConfiguration: []gads.BiddingStrategyConfiguration{
					gads.BiddingStrategyConfiguration{
						StrategyType: "MANUAL_CPC",
						Bids: []gads.Bid{
							gads.Bid{
								Type:   "CpcBid",
								Amount: 2000000,
							},
						},
					},
				},
			},
		},
	}*/

	ago := gads.AdGroupCriterionOperations{
		"SET": {
			gads.BiddableAdGroupCriterion{
				AdGroupId: 1234567890,
				Criterion: gads.KeywordCriterion{
					Id: 1234567890,
				},
				BiddingStrategyConfiguration: &gads.BiddingStrategyConfiguration{
					Bids: []gads.Bid{
						gads.Bid{
							Type:   "CpcBid",
							Amount: 3000000,
						},
					},
				},
			},
			gads.BiddableAdGroupCriterion{
				AdGroupId: 1234567890,
				Criterion: gads.KeywordCriterion{
					Id: 1234567890,
				},
				BiddingStrategyConfiguration: &gads.BiddingStrategyConfiguration{
					Bids: []gads.Bid{
						gads.Bid{
							Type:   "CpcBid",
							Amount: 2000000,
						},
					},
				},
			},
		},
	}

	if resp, err := ms.Mutate(ago, policy); err == nil {
		jobId := resp.Id

		// loop
		for {
			// recheck every 5 seconds
			time.Sleep(5 * time.Second)

			jobSelector := gads.BulkMutateJobSelector{JobIds: []int64{jobId}}

			result, err := ms.Get(jobSelector)

			if err != nil {
				panic(err)
			}

			if result.Status == "COMPLETED" {
				break
			}

			if result.Status == "FAILED" {
				// probably do something else here
				panic("Job result failed")
			}
		}

		jobSelector := gads.BulkMutateJobSelector{JobIds: []int64{jobId}}
		result, err := ms.GetResult(jobSelector)

		if err != nil {
			panic(err)
		}
		jsonResult, _ := json.Marshal(result)
		fmt.Println(result)
		fmt.Println(string(jsonResult))

	} else {
		// handle err
		panic(err)
	}

}

func rand_str(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

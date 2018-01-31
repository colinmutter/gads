package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	gads "github.com/Getsidecar/gads/v201710"
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

	// Batch Job
	bs := gads.NewBatchJobService(&config.Auth)

	/*ago := gads.AdGroupCriterionOperations{
	"SET":	gads.AdGroupCriterions {
				gads.BiddableAdGroupCriterion{
					AdGroupId: 11400713462,
					Criterion: gads.KeywordCriterion{
						Id: 47652524966,
					},
					BiddingStrategyConfiguration: &gads.BiddingStrategyConfiguration{
						Bids: []gads.Bid{
							gads.Bid{
								Type:   "CpcBid",
								Amount: 5372200,
							},
						},
					},
					UrlCustomParameters: gads.CustomParameters{
						CustomParameters: []gads.CustomParameter{
							gads.CustomParameter{
								Key:      "foo",
								Value:    "bar",
								IsRemove: false,
							},
						},
						DoReplace: false,
					},
				},
				gads.BiddableAdGroupCriterion{
					AdGroupId: 11400713462,
					Criterion: gads.KeywordCriterion{
						Id: 47652524366,
					},
					BiddingStrategyConfiguration: &gads.BiddingStrategyConfiguration{
						Bids: []gads.Bid{
							gads.Bid{
								Type:   "CpcBid",
								Amount: 5372200,
							},
						},
					},
					UrlCustomParameters: gads.CustomParameters{
						CustomParameters: []gads.CustomParameter{
							gads.CustomParameter{
								Key:      "foo",
								Value:    "bar",
								IsRemove: false,
							},
						},
						DoReplace: false,
					},
				},
			},
	}*/
	// Creating AdGroups
	ago := gads.AdGroupOperations{
		"ADD": {},
	}

	// stress test uploading
	var adGroupNum int = 10000

	for i := 0; i < adGroupNum; i++ {
		ago["ADD"] = append(ago["ADD"], gads.AdGroup{
			Name:       "test ad group " + rand_str(10),
			Status:     "PAUSED",
			CampaignId: campaignId,
		})
	}

	var operations []interface{}
	operations = append(operations, ago)

	bjo := gads.BatchJobOperations{
		BatchJobOperations: []gads.BatchJobOperation{
			gads.BatchJobOperation{
				Operator: "ADD",
				Operand:  gads.BatchJob{},
			},
		},
	}

	if resp, err := bs.Mutate(bjo); err == nil {

		bjh := gads.NewBatchJobHelper(&config.Auth)
		err = bjh.UploadBatchJobOperations(operations, *resp[0].UploadUrl)

		if err != nil {
			panic(err)
		}

		jobId := resp[0].Id
		batchJobs := gads.BatchJobPage{}

		// loop
		for {
			// recheck every 5 seconds
			time.Sleep(5 * time.Second)
			selector := gads.Selector{
				Fields: []string{
					"Id",
					"Status",
					"DownloadUrl",
					"ProcessingErrors",
					"ProgressStats",
				},
				Predicates: []gads.Predicate{
					{"Id", "EQUALS", []string{strconv.FormatInt(jobId, 10)}},
				},
			}

			// more than likely you'll want to have some logic to loop through these if you have multiple batch jobs, but since only one we just want to grab the first one
			batchJobs, err = bs.Get(selector)

			if err != nil {
				panic(err)
			}

			if batchJobs.BatchJobs[0].Status == "DONE" {
				break
			} else if batchJobs.BatchJobs[0].Status == "CANCELED" {
				panic("Job was canceled")
			}
		}

		if batchJobs.BatchJobs[0].DownloadUrl.Url != "" {
			// get the job
			mutateResult, err := bjh.DownloadBatchJob(*batchJobs.BatchJobs[0].DownloadUrl)

			if err != nil {
				panic(err)
			}

			fmt.Println(mutateResult)
			jsonResult, _ := json.Marshal(mutateResult)

			fmt.Println(string(jsonResult))
		}
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

package main

import (
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

	config.Auth.PartialFailure = true

	service := gads.NewAdGroupService(&config.Auth)

	adGroupList := []gads.AdGroupLabel{
		gads.AdGroupLabel{
			AdGroupId: 1234567890,
			LabelId:   1234567,
		},
		gads.AdGroupLabel{
			AdGroupId: 1234567891,
			LabelId:   1234567,
		},
		gads.AdGroupLabel{
			AdGroupId: 1234567892,
			LabelId:   1234567,
		},
	}

	adGroupLabels, err := service.MutateLabel(
		gads.AdGroupLabelOperations{
			"ADD": adGroupList,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(adGroupLabels)

}

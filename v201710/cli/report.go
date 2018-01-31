package main

import (
	"flag"
	"log"

	gads "github.com/Getsidecar/gads/v201710"
)

var configJson = flag.String("oauth", "./oauth.json", "API credentials")

// The query you want to run
var awql string = "Select AdGroupId,Id,CreativeQualityScore,PostClickQualityScore,SearchPredictedCtr,QualityScore FROM KEYWORDS_PERFORMANCE_REPORT DURING YESTERDAY"

func main() {
	flag.Parse()
	config, err := gads.NewCredentialsFromFile(*configJson)
	if err != nil {
		log.Fatal(err)
	}

	// Report Service
	rs := gads.NewReportDownloadService(&config.Auth)

	res, err := rs.AWQL(awql, "CSV")

	if err != nil {
		log.Panicln(err)
	}

	log.Println(res)
}

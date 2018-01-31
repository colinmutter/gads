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

	tis := gads.NewTargetingIdeaService(&config.Auth)

	ads, _, err := tis.Get(
		gads.TargetingIdeaSelector{
			IdeaType:    "KEYWORD",
			RequestType: "IDEAS",
			RequestedAttributeTypes: []string{
				"CATEGORY_PRODUCTS_AND_SERVICES",
				"COMPETITION",
				"EXTRACTED_FROM_WEBPAGE",
				"IDEA_TYPE",
				"KEYWORD_TEXT",
				"SEARCH_VOLUME",
				"AVERAGE_CPC",
				"TARGETED_MONTHLY_SEARCHES",
			},
			SearchParameters: []gads.SearchParameter{
				gads.CategoryProductsAndServicesSearchParameter{
					CategoryID: 51,
				},
				gads.CompetitionSearchParameter{
					Levels: []string{"MEDIUM", "HIGH"},
				},
				gads.IdeaTextFilterSearchParameter{
					Included: []string{},
					Excluded: []string{"red herring e3a2b4b7", "red herring 418f2d72"},
				},
				gads.IncludeAdultContentSearchParameter{},
				gads.LanguageSearchParameter{
					Languages: []gads.LanguageCriterion{
						gads.LanguageCriterion{Id: 1000},
					},
				},
				gads.LocationSearchParameter{
					Locations: []gads.Location{
						gads.Location{Id: 2840},
						gads.Location{Id: 2124},
					},
				},
				gads.NetworkSearchParameter{
					NetworkSetting: gads.NetworkSetting{
						TargetGoogleSearch: true,
					},
				},
				gads.RelatedToQuerySearchParameter{
					Queries: []string{"blue herring", "test"},
				},
				gads.RelatedToUrlSearchParameter{
					Urls: []string{"https://google.com"},
				},
				gads.SearchVolumeSearchParameter{
					Minimum: 3900000,
					Maximum: 5445394,
				},
				// gads.SeedAdGroupIdSearchParameter{
				// 	AdGroupID: 123456,
				// },
			},
			Paging: gads.Paging{
				Offset: 0,
				Limit:  3,
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

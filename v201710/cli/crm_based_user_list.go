package main

import (
	"flag"
	"fmt"
	"log"

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

	// User List Service
	auls := gads.NewAdwordsUserListService(&config.Auth)

	crmList := gads.NewCrmBasedUserList("Test List", "Just a list to test with", 0, "http://mytest.com/optout")

	ops := gads.UserListOperations{
		Operations: []gads.Operation{
			gads.Operation{
				Operator: "ADD",
				Operand:  crmList,
			},
		},
	}

	resp, err := auls.Mutate(ops)

	if err != nil {
		panic(err)
	}

	fmt.Println(resp[0].Id)

	mmo := gads.NewMutateMembersOperand()
	mmo.UserListId = resp[0].Id

	var members []string
	members = append(members, "brian@test.com")
	members = append(members, "test@test.com")

	mmo.Members = members

	mutateMembersOperations := gads.MutateMembersOperations{
		Operations: []gads.Operation{
			gads.Operation{
				Operator: "ADD",
				Operand:  mmo,
			},
		},
	}

	lists, err := auls.MutateMembers(mutateMembersOperations)

	if err != nil {
		panic(err)
	}

	fmt.Println(lists)
}

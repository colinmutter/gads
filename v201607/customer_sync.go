package v201607

import "encoding/xml"

type (
	// CustomerSyncService struct implements the Auth struct
	CustomerSyncService struct {
		Auth
	}

	// AdGroupChangeData Holds information about a changed adgroup
	AdGroupChangeData struct {
		AdGroupID                         int64   `xml:"adGroupId"`
		AdGroupChangeStatus               string  `xml:"adGroupChangeStatus"`
		ChangedAds                        []int64 `xml:"changedAds"`
		ChangedCriteria                   []int64 `xml:"changedCriteria"`
		RemovedCriteria                   []int64 `xml:"removedCriteria"`
		ChangedFeeds                      []int64 `xml:"changedFeeds"`
		ChangedAdGroupBidModifierCriteria []int64 `xml:"changedAdGroupBidModifierCriteria"`
		RemovedAdGroupBidModifierCriteria []int64 `xml:"removedAdGroupBidModifierCriteria"`
	}

	// CampaignChangeData Holds information about a changed campaign and any ad groups under that have changed.
	CampaignChangeData struct {
		CampaignID              int64               `xml:"campaignId"`
		CampaignChangeStatus    string              `xml:"campaignChangeStatus"`
		ChangedAdGroups         []AdGroupChangeData `xml:"changedAdGroups"`
		AddedCampaignCriteria   []int64             `xml:"addedCampaignCriteria"`
		RemovedCampaignCriteria []int64             `xml:"removedCampaignCriteria"`
		ChangedFeeds            []int64             `xml:"changedFeeds"`
		RemovedFeeds            []int64             `xml:"removedFeeds"`
	}

	// FeedChangeData Holds information about a changed feed and any feeds items within the feed.
	FeedChangeData struct {
		FeedID           int64  `xml:"feedId"`
		FeedChangeStatus string `xml:"feedChangeStatus"`
		ChangedFeedItems int64  `xml:"changedFeedItems"`
		RemovedFeedItems int64  `xml:"removedFeedItems"`
	}

	// CustomerChangeData Holds information about changes to a customer
	CustomerChangeData struct {
		ChangedCampaigns    []CampaignChangeData `xml:"changedCampaigns"`
		ChangedFeeds        []FeedChangeData     `xml:"changedFeeds"`
		LastChangeTimestamp string               `xml:"lastChangeTimestamp"`
	}
)

// NewCustomerSyncService implements a new CustomerSyncService struct
func NewCustomerSyncService(auth *Auth) *CustomerSyncService {
	return &CustomerSyncService{Auth: *auth}
}

// Get Returns information about changed entities inside a customer's account.
func (s *CustomerSyncService) Get(selector Selector) (changes []CustomerChangeData, err error) {
	respBody, err := s.Auth.request(
		customerSyncServiceUrl,
		"get",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "get",
			},
			Sel: selector,
		},
	)
	if err != nil {
		return changes, err
	}
	getResp := struct {
		Changes []CustomerChangeData `xml:"rval"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return changes, err
	}
	return getResp.Changes, err
}

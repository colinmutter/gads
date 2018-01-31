package v201710

import "encoding/xml"

type CampaignSharedSetService struct {
	Auth
}

func NewCampaignSharedSetService(auth *Auth) *CampaignSharedSetService {
	return &CampaignSharedSetService{Auth: *auth}
}

type CampaignSharedSet struct {
	SharedSetId   int64  `xml:"sharedSetId,omitempty"`
	CampaignId    int64  `xml:"campaignId"`
	SharedSetName string `xml:"sharedSetName"`
	SharedSetType string `xml:"sharedSetType"`
	CampaignName  string `xml:"campaignName"`
	Status        string `xml:"status"`
}

type CampaignSharedSetOperation struct {
	Operator string            `xml:"operator,omitempty"`
	Operand  CampaignSharedSet `xml:"operand,omitempty"`
}

func (s CampaignSharedSetService) Get(selector Selector) (sharedSets []CampaignSharedSet, totalCount int64, err error) {
	selector.XMLName = xml.Name{baseUrl, "selector"}
	respBody, err := s.Auth.request(
		campaignSharedSetServiceUrl,
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
		return sharedSets, totalCount, err
	}
	getResp := struct {
		Size       int64               `xml:"rval>totalNumEntries"`
		SharedSets []CampaignSharedSet `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return sharedSets, totalCount, err
	}
	return getResp.SharedSets, getResp.Size, err
}

func (s CampaignSharedSetService) Mutate(operations []CampaignSharedSetOperation) error {
	mutateRequest := struct {
		XMLName xml.Name
		Ops     []CampaignSharedSetOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: operations}
	_, err := s.Auth.request(campaignSharedSetServiceUrl, "mutate", mutateRequest)
	return err
}

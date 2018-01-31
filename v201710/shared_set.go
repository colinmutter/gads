package v201710

import "encoding/xml"

type SharedSetService struct {
	Auth
}

func NewSharedSetService(auth *Auth) *SharedSetService {
	return &SharedSetService{Auth: *auth}
}

type SharedSet struct {
	Id             int64  `xml:"sharedSetId,omitempty"`
	Name           string `xml:"name,omitempty"`
	Type           string `xml:"type,omitempty"`
	MemberCount    int    `xml:"memberCount,omitempty"`
	ReferenceCount int    `xml:"referenceCount,omitempty"`
	Status         string `xml:"status,omitempty"`
}

type SharedSetOperation struct {
	Operator string    `xml:"operator,omitempty"`
	Operand  SharedSet `xml:"operand,omitempty"`
}

func (s SharedSetService) Get(selector Selector) (sharedSets []SharedSet, totalCount int64, err error) {
	selector.XMLName = xml.Name{baseUrl, "selector"}
	respBody, err := s.Auth.request(
		sharedSetServiceUrl,
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
		Size       int64       `xml:"rval>totalNumEntries"`
		SharedSets []SharedSet `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return sharedSets, totalCount, err
	}
	return getResp.SharedSets, getResp.Size, err
}

func (s SharedSetService) Mutate(operations []SharedSetOperation) ([]SharedSet, error) {
	mutateRequest := struct {
		XMLName xml.Name
		Ops     []SharedSetOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: operations}

	respBody, err := s.Auth.request(sharedSetServiceUrl, "mutate", mutateRequest)

	if err != nil {
		return nil, err
	}

	getResp := struct {
		ListReturnValueType string      `xml:"rval>ListReturnValue.Type"`
		SharedSets          []SharedSet `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return nil, err
	}
	return getResp.SharedSets, nil
}

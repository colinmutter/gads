package v201710

import (
	"encoding/xml"
	"fmt"
)

type SharedCriterionService struct {
	Auth
}

func NewSharedCriterionService(auth *Auth) *SharedCriterionService {
	return &SharedCriterionService{Auth: *auth}
}

type SharedCriterion struct {
	SharedSetId int64     `xml:"sharedSetId,omitempty"`
	Negative    bool      `xml:"negative,omitempty"`
	Criterion   Criterion `xml:"criterion,omitempty"`
}

type SharedCriterionOperation struct {
	Operator string          `xml:"operator,omitempty"`
	Operand  SharedCriterion `xml:"operand,omitempty"`
}

func (s SharedCriterionService) Get(selector Selector) (sharedCriteria []SharedCriterion, totalCount int64, err error) {
	selector.XMLName = xml.Name{baseUrl, "selector"}
	respBody, err := s.Auth.request(
		sharedCriterionServiceUrl,
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
		return sharedCriteria, totalCount, err
	}
	getResp := struct {
		Size           int64             `xml:"rval>totalNumEntries"`
		SharedCriteria []SharedCriterion `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return sharedCriteria, totalCount, err
	}
	return getResp.SharedCriteria, getResp.Size, err
}

func (s SharedCriterionService) Mutate(operations []SharedCriterionOperation) error {
	mutateRequest := struct {
		XMLName xml.Name
		Ops     []SharedCriterionOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: operations}
	_, err := s.Auth.request(sharedCriterionServiceUrl, "mutate", mutateRequest)
	return err
}

func (s *SharedCriterion) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			tag := start.Name.Local
			switch tag {
			case "sharedSetId":
				if err := dec.DecodeElement(&s.SharedSetId, &start); err != nil {
					return err
				}
			case "negative":
				if err := dec.DecodeElement(&s.Negative, &start); err != nil {
					return err
				}
			case "criterion":
				criterion, err := criterionUnmarshalXML(dec, start)
				if err != nil {
					return err
				}
				s.Criterion = criterion
			default:
				return fmt.Errorf("unknown BiddableAdGroupCriterion field %s", tag)
			}
		}
	}
	return nil
}

func (s SharedCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(
		start.Attr,
		xml.Attr{
			xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"},
			"SharedCriterion",
		},
	)
	e.EncodeToken(start)
	e.EncodeElement(&s.SharedSetId, xml.StartElement{Name: xml.Name{baseUrl, "sharedSetId"}})
	criterionMarshalXML(s.Criterion, e)
	e.EncodeElement(&s.Negative, xml.StartElement{Name: xml.Name{baseUrl, "negative"}})
	e.EncodeToken(start.End())
	return nil
}

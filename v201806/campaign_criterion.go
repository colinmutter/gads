package v201806

import (
	"encoding/xml"
	"fmt"
)

type CampaignCriterionService struct {
	Auth
}

func NewCampaignCriterionService(auth *Auth) *CampaignCriterionService {
	return &CampaignCriterionService{Auth: *auth}
}

type CampaignCriterion struct {
	CampaignId  int64     `xml:"campaignId"`
	IsNegative  bool      `xml:"isNegative,omitempty"`
	Criterion   Criterion `xml:"criterion"`
	BidModifier float64   `xml:"bidModifier,omitempty"`
	Errors      []error   `xml:"-"`
	Type        string    `xml:"-"`
	Id          int64     `xml:"-"`
}

type NegativeCampaignCriterion CampaignCriterion

func (cc CampaignCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	//fmt.Printf("processing -> %#v\n",ncc)
	start.Attr = append(
		start.Attr,
		xml.Attr{
			xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"},
			"CampaignCriterion",
		},
	)
	e.EncodeToken(start)
	e.EncodeElement(&cc.CampaignId, xml.StartElement{Name: xml.Name{"", "campaignId"}})

	if cc.Criterion == nil {
		var ok bool
		if cc.Criterion, ok = CriterionFromIdAndType(cc.Id, cc.Type); !ok {
			return fmt.Errorf("missing criterion")
		}
	}
	if err := criterionMarshalXML(cc.Criterion, e); err != nil {
		return err
	}
	if cc.BidModifier != 0 {
		e.EncodeElement(&cc.BidModifier, xml.StartElement{Name: xml.Name{"", "bidModifier"}})
	}

	e.EncodeToken(start.End())
	return nil
}

type CampaignCriterions []interface{}
type CampaignCriterionOperations map[string]CampaignCriterions

func (ncc NegativeCampaignCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	//fmt.Printf("processing -> %#v\n",ncc)
	start.Attr = append(
		start.Attr,
		xml.Attr{
			xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"},
			"NegativeCampaignCriterion",
		},
	)
	e.EncodeToken(start)
	e.EncodeElement(&ncc.CampaignId, xml.StartElement{Name: xml.Name{"", "campaignId"}})
	criterionMarshalXML(ncc.Criterion, e)
	e.EncodeToken(start.End())
	return nil
}

func (ccs *CampaignCriterions) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	cc := CampaignCriterion{}

	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "campaignId":
				if err := dec.DecodeElement(&cc.CampaignId, &start); err != nil {
					return err
				}
			case "criterion":
				criterion, err := criterionUnmarshalXML(dec, start)
				if err != nil {
					return err
				}
				cc.Id, cc.Type, _ = CriterionIdAndType(criterion)
				cc.Criterion = criterion
			case "bidModifier":
				if err := dec.DecodeElement(&cc.BidModifier, &start); err != nil {
					return err
				}
			case "isNegative":
				if err := dec.DecodeElement(&cc.IsNegative, &start); err != nil {
					return err
				}
			}
		}
	}
	if cc.IsNegative {
		*ccs = append(*ccs, NegativeCampaignCriterion(cc))
	} else {
		*ccs = append(*ccs, cc)
	}
	return nil
}

/*
func NewNegativeCampaignCriterion(campaignId int64, bidModifier float64, criterion interface{}) CampaignCriterion {
  return CampaignCriterion{
    CampaignId: campaignId,
    Criterion: criterion,
    BidModifier: bidModifier
  }
  switch c := criterion.(type) {
  case AdScheduleCriterion:
  case AgeRangeCriterion:
  case ContentLabelCriterion:
  case GenderCriterion:
  case KeywordCriterion:
  case LanguageCriterion:
  case LocationCriterion:
  case MobileAppCategoryCriterion:
  case MobileApplicationCriterion:
  case MobileDeviceCriterion:
  case OperatingSystemVersionCriterion:
  case PlacementCriterion:
  case PlatformCriterion:
  case ProductCriterion:
  case ProximityCriterion:
  case UserInterestCriterion:
    cc.Criterion = criterion
  case UserListCriterion:
    cc.Criterion = criterion
  case VerticalCriterion:
  }
}
*/

func (s *CampaignCriterionService) Get(selector Selector) (campaignCriterions CampaignCriterions, totalCount int64, err error) {
	selector.XMLName = xml.Name{baseUrl, "serviceSelector"}
	getResp := struct {
		XMLName            xml.Name
		Size               int64              `xml:"rval>totalNumEntries"`
		CampaignCriterions CampaignCriterions `xml:"rval>entries"`
	}{}

	err = s.Auth.do(
		campaignCriterionServiceUrl,
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
		&getResp,
	)
	if err != nil {
		return campaignCriterions, totalCount, err
	}
	return getResp.CampaignCriterions, getResp.Size, err
}

type CampaignCriterionOperation struct {
	Action            string      `xml:"operator"`
	CampaignCriterion interface{} `xml:"operand"`
}

func (s *CampaignCriterionService) MutateOperations(operations []CampaignCriterionOperation) (CampaignCriterions, error) {
	mutation := struct {
		XMLName xml.Name
		Ops     []CampaignCriterionOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: operations,
	}

	mutateResp := struct {
		XMLName            xml.Name
		CampaignCriterions CampaignCriterions `xml:"rval>value"`
	}{}
	err := s.Auth.do(campaignCriterionServiceUrl, "mutate", mutation, &mutateResp)
	if err != nil {
		/*
			    switch t := err.(type) {
			    case *ErrorsType:
				    for action, campaignCriterions := range campaignCriterionOperations {
					    for _, campaignCriterion := range campaignCriterions {
			          campaignCriterions = append(campaignCriterions,campaignCriterion)
			        }
			      }
			      for _, aef := range t.ApiExceptionFaults {
			        for _,e := range aef.Errors {
			          switch et := e.(type) {
			          case CriterionError:
			            offset, err := strconv.ParseInt(strings.Trim(et.FieldPath,"abcdefghijklmnop.]["),10,64)
			            if err != nil {
			              return CampaignCriterions{}, err
			            }
			            cc := campaignCriterions[offset]
			            switch c := cc.(type) {
			            case CampaignCriterion:
			              CampaignCriterion(campaignCriterions[offset]).Errors = append(campaignCriterions[offset].(CampaignCriterion).Errors,fmt.Errorf(et.Reason))
			            case NegativeCampaignCriterion:
			              NegativeCampaignCriterion(campaignCriterions[offset]).Errors = append(NegativeCampaignCriterion(campaignCriterions[offset].Errors),fmt.Errorf(et.Reason))
			            }
			          }
			        }
			      }
			    default:
		*/
		return nil, err
		//}
	}
	return mutateResp.CampaignCriterions, err
}

func (s *CampaignCriterionService) Mutate(campaignCriterionOperations CampaignCriterionOperations) (campaignCriterions CampaignCriterions, err error) {
	operations := []CampaignCriterionOperation{}
	for action, campaignCriterions := range campaignCriterionOperations {
		for _, campaignCriterion := range campaignCriterions {
			operations = append(operations,
				CampaignCriterionOperation{
					Action:            action,
					CampaignCriterion: campaignCriterion,
				},
			)
		}
	}

	return s.MutateOperations(operations)
}

func (s *CampaignCriterionService) Query(query string) (campaignCriterions CampaignCriterions, totalCount int64, err error) {
	respBody, err := s.Auth.request(
		campaignCriterionServiceUrl,
		"query",
		AWQLQuery{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "query",
			},
			Query: query,
		},
	)

	if err != nil {
		return campaignCriterions, totalCount, err
	}

	getResp := struct {
		Size               int64              `xml:"rval>totalNumEntries"`
		CampaignCriterions CampaignCriterions `xml:"rval>entries"`
	}{}

	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return campaignCriterions, totalCount, err
	}
	return getResp.CampaignCriterions, getResp.Size, err
}

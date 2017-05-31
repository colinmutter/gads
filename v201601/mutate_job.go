package v201601

import (
	"encoding/xml"
	"errors"
	"reflect"
	"strings"
)

type MutateJobService struct {
	Auth
}

type BulkMutateJobPolicy struct {
	PrerequisiteJobIds []int64 `xml:"prerequisiteJobIds"`
}

type Operation struct {
	Operator string      `xml:"https://adwords.google.com/api/adwords/cm/v201601 operator"`
	Operand  interface{} `xml:"operand"`
	Xsi_type string      `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr,omitempty"`
}

type SimpleMutateJob struct {
	Id     int64  `xml:"rval>id" json:",string"`
	Status string `xml:"rval>status"`
}

type BulkMutateJobSelector struct {
	JobIds          []int64 `xml:"jobIds"`
	ResultPartIndex int     `xml:"ResultPartIndex,omitempty"`
	Xsi_type        string  `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr,omitempty"`
}

type SimpleMutateResponse struct {
	Results []Result      `xml:"rval>SimpleMutateResult>results"`
	Errors  []EntityError `xml:"rval>SimpleMutateResult>errors"`
}

type Result struct {
	AdGroup               *AdGroup               `xml:"AdGroup,omitempty"`
	AdGroupAdLabel        *AdGroupAdLabel        `xml:"AdGroupAdLabel,omitempty"`
	AdGroupCriterion      *AdGroupCriterions     `xml:"AdGroupCriterion,omitempty"`
	AdGroupCriterionLabel *AdGroupCriterionLabel `xml:"AdGroupCriterionLabel,omitempty"`
	AdGroupLabel          *AdGroupLabel          `xml:"AdGroupLabel,omitempty"`
	Budget                *Budget                `xml:"Budget,omitempty"`
	CampaignCriterion     *CampaignCriterion     `xml:"CampaignCriterion,omitempty"`
	CampaignLabel         *CampaignLabel         `xml:"CampaignLabel,omitempty"`
	Campaign              *Campaign              `xml:"Campaign,omitempty"`
	Label                 *Label                 `xml:"Label,omitempty"`
	Media                 *Media                 `xml:"Media,omitempty"`
}

func NewMutateJobService(auth *Auth) *MutateJobService {
	return &MutateJobService{Auth: *auth}
}

// Mutate bulk mutates the operations passed in
// Note: this would be a lot easier if there was an operation struct we were passing in
// But that would require refactoring everything and probably breaking dependencies... perhaps another day
func (s *MutateJobService) Mutate(jobOperations interface{}, policy *BulkMutateJobPolicy) (mutateResp SimpleMutateJob, err error) {

	var operations []Operation

	if operationType, valid := getXsiType(reflect.ValueOf(jobOperations).Type().String()); valid {
		switch reflect.TypeOf(jobOperations).Kind() {
		case reflect.Map:
			ops := reflect.ValueOf(jobOperations)

			keys := ops.MapKeys()

			for _, action := range keys {
				jobs := ops.MapIndex(action)

				for i := 0; i < jobs.Len(); i++ {

					operations = append(operations,
						Operation{
							Operator: action.String(),
							Operand:  jobs.Index(i).Interface(),
							Xsi_type: operationType,
						},
					)
				}
			}
		}

		// make sure we actually have operations to send
		if len(operations) > 0 {
			mutation := struct {
				XMLName xml.Name
				Ops     []Operation          `xml:"operations"`
				Policy  *BulkMutateJobPolicy `xml:"policy"`
			}{
				XMLName: xml.Name{
					Space: baseUrl,
					Local: "mutate",
				},
				Ops:    operations,
				Policy: policy,
			}

			respBody, err := s.Auth.request(mutateJobServiceUrl, "mutate", mutation)

			if err != nil {
				return mutateResp, err
			}

			err = xml.Unmarshal([]byte(respBody), &mutateResp)
		}

	} else {
		err = errors.New("Invalid Operation type passed in")
	}

	return mutateResp, err
}

// Get queries mutation results of existing jobs
func (s *MutateJobService) Get(jobSelector BulkMutateJobSelector) (mutateResp SimpleMutateJob, err error) {

	jobSelector.Xsi_type = "BulkMutateJobSelector"

	respBody, err := s.Auth.request(
		mutateJobServiceUrl,
		"get",
		struct {
			XMLName xml.Name
			Sel     BulkMutateJobSelector `xml:"selector"`
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "get",
			},
			Sel: jobSelector,
		},
	)

	if err != nil {
		return mutateResp, err
	}

	err = xml.Unmarshal([]byte(respBody), &mutateResp)

	return mutateResp, err
}

// Get mutation results of a completed job
func (s *MutateJobService) GetResult(jobSelector BulkMutateJobSelector) (mutateResp SimpleMutateResponse, err error) {

	jobSelector.Xsi_type = "BulkMutateJobSelector"

	respBody, err := s.Auth.request(
		mutateJobServiceUrl,
		"getResult",
		struct {
			XMLName xml.Name
			Sel     BulkMutateJobSelector `xml:"selector"`
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "getResult",
			},
			Sel: jobSelector,
		},
	)

	if err != nil {
		return mutateResp, err
	}

	err = xml.Unmarshal([]byte(respBody), &mutateResp)

	return mutateResp, err
}

// getXsiType validates the schema instance type and returns it since Bulk Mutate requires it to be set
func getXsiType(objectName string) (string, bool) {
	switch {
	case strings.Contains(objectName, "AdGroupAdLabelOperation"):
		return "AdGroupAdLabelOperation", true
	case strings.Contains(objectName, "AdGroupAdOperation"):
		return "AdGroupAdOperation", true
	case strings.Contains(objectName, "AdGroupBidModifierOperation"):
		return "AdGroupBidModifierOperation", true
	case strings.Contains(objectName, "AdGroupCriterionLabelOperation"):
		return "AdGroupCriterionLabelOperation", true
	case strings.Contains(objectName, "AdGroupCriterionOperation"):
		return "AdGroupCriterionOperation", true
	case strings.Contains(objectName, "AdGroupLabelOperation"):
		return "AdGroupLabelOperation", true
	case strings.Contains(objectName, "AdGroupOperation"):
		return "AdGroupOperation", true
	case strings.Contains(objectName, "BudgetOperation"):
		return "BudgetOperation", true
	case strings.Contains(objectName, "CampaignAdExtensionOperation"):
		return "CampaignAdExtensionOperation", true
	case strings.Contains(objectName, "CampaignCriterionOperation"):
		return "CampaignCriterionOperation", true
	case strings.Contains(objectName, "CampaignLabelOperation"):
		return "CampaignLabelOperation", true
	case strings.Contains(objectName, "CampaignOperation"):
		return "CampaignOperation", true
	case strings.Contains(objectName, "FeedItemOperation"):
		return "FeedItemOperation", true
	default:
		return "", false
	}
}

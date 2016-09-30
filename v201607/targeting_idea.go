package v201607

import "encoding/xml"

type (
	// TargetIdeaService struct implements the Auth struct
	TargetIdeaService struct {
		Auth
	}

	// AttributeType This represents an entry in a map with a key of type Type and value of type Attribute.
	AttributeType struct {
		Key   string `xml:"key"`
		Value string `xml:"value"`
	}

	// TargetingIdea Represents a TargetingIdea returned by search criteria
	// specified in the TargetingIdeaSelector.
	// Targeting ideas are keywords or placements that are similar to those the user inputs.
	TargetingIdea struct {
		Data []AttributeType `xml:"data"`
	}

	// TargetingIdeaPage contains a subset of TargetingIdeas
	// from the search criteria specified by a TargetingIdeaSelector
	TargetingIdeaPage struct {
		TotalNumEntries int             `xml:"totalNumEntries"`
		Entries         []TargetingIdea `xml:"entries"`
	}
)

// NewTargetIdeaService initializes a new TargetIdeaService struct
func NewTargetIdeaService(auth *Auth) *TargetIdeaService {
	return &TargetIdeaService{Auth: *auth}
}

// Get Returns a page of ideas that match the query described by the specified TargetingIdeaSelector.
// The selector must specify a paging value, with numberResults set to 800 or less.
// Large result sets must be composed through multiple calls to this method,
// advancing the paging startIndex value by numberResults with each call.
func (s *TargetIdeaService) Get(selector Selector) (ideas []TargetingIdeaPage, err error) {
	respBody, err := s.Auth.request(
		targetingIdeaServiceUrl,
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
		return ideas, err
	}
	getResp := struct {
		Ideas []TargetingIdeaPage `xml:"rval"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return ideas, err
	}
	return getResp.Ideas, err
}

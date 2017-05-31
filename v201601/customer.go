package v201601

import (
	"encoding/xml"
)

type CustomerService struct {
	Auth
}

type Customer struct {
	CustomerId                 int64                      `xml:"customerId"`
	CurrencyCode               string                     `xml:"currencyCode"`
	DateTimeZone               string                     `xml:"dateTimeZone"`
	DescriptiveName            string                     `xml:"descriptiveName"`
	CompanyName                string                     `xml:"companyName"`
	CanManageClients           bool                       `xml:"canManageClients"`
	TestAccount                bool                       `xml:"testAccount"`
	AutoTaggingEnabled         bool                       `xml:"autoTaggingEnabled"`
	TrackingUrlTemplate        string                     `xml:"trackingUrlTemplate,omitempty"`
	ConversionTrackingSettings ConversionTrackingSettings `xml:"conversionTrackingSettings"`
	RemarketingSettings        RemarketingSettings        `xml:"remarketingSettings"`
}

type RemarketingSettings struct {
	Snippet string `xml:"snippet"`
}

func NewCustomerService(auth *Auth) *CustomerService {
	return &CustomerService{Auth: *auth}
}

func (s *CustomerService) Get(selector Selector) (customer Customer, err error) {
	selector.XMLName = xml.Name{baseMcmUrl, "serviceSelector"}
	respBody, err := s.Auth.request(
		customerServiceUrl,
		"get",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{
				Space: baseMcmUrl,
				Local: "get",
			},
			Sel: selector,
		},
	)
	if err != nil {
		return customer, err
	}
	getResp := struct {
		Customer Customer `xml:"rval"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return customer, err
	}
	return getResp.Customer, nil
}

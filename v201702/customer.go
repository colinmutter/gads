package v201609

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

func (s *CustomerService) GetCustomers() (customers []Customer, err error) {
	respBody, err := s.Auth.request(
		customerServiceUrl,
		"getCustomers",
		struct {
			XMLName xml.Name
		}{
			XMLName: xml.Name{
				Space: baseMcmUrl,
				Local: "getCustomers",
			},
		},
	)
	if err != nil {
		return customers, err
	}
	getResp := struct {
		Customers []Customer `xml:"rval"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return customers, err
	}
	return getResp.Customers, nil
}

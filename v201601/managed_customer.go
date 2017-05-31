package v201601

import (
	"encoding/xml"
)

type ManagedCustomer struct {
	Name                  string         `xml:"name"`
	CompanyName           string         `xml:"companyName"`
	CustomerId            int64          `xml:"customerId"`
	CanManageClients      bool           `xml:"canManageClients"`
	CurrencyCode          string         `xml:"currencyCode"`
	DateTimeZone          string         `xml:"dateTimeZone"`
	TestAccount           bool           `xml:"testAccount"`
	AccountLabels         []AccountLabel `xml:"accountLabels"`
	ExcludeHiddenAccounts bool           `xml:"excludeHiddenAccounts"`
}

type ManagedCustomerLink struct {
	ManagerCustomerId      int64  `xml:"managerCustomerId"`
	ClientCustomerId       int64  `xml:"clientCustomerId"`
	LinkStatus             string `xml:"linkStatus"`
	PendingDescriptiveName string `xml:"pendingDescriptiveName"`
	IsHidden               bool   `xml:isHidden"`
}

type ManagedCustomerPage struct {
	Size                 int64                 `xml:"rval>totalNumEntries"`
	ManagedCustomers     []ManagedCustomer     `xml:"rval>entries"`
	ManagedCustomerLinks []ManagedCustomerLink `xml:"rval>links"`
}

type AccountLabel struct {
	Id   int64  `xml:"id"`
	Name string `xml:"name"`
}

type ManagedCustomerService struct {
	Auth
}

func NewManagedCustomerService(auth *Auth) *ManagedCustomerService {
	return &ManagedCustomerService{Auth: *auth}
}

func (s *ManagedCustomerService) Get(selector Selector) (managedCustomerPage ManagedCustomerPage, totalCount int64, err error) {
	selector.XMLName = xml.Name{baseMcmUrl, "serviceSelector"}
	respBody, err := s.Auth.request(
		managedCustomerServiceUrl,
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
		return managedCustomerPage, totalCount, err
	}
	getResp := ManagedCustomerPage{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return managedCustomerPage, totalCount, err
	}
	return getResp, totalCount, nil
}

package v201710

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type TestClient struct {
	res *http.Response
}

func (s *TestClient) Do(req *http.Request) (*http.Response, error) {
	return s.res, nil
}

func TestReportDownloadAuthError(t *testing.T) {
	body := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?> <reportDownloadError> <ApiError> <type>AuthorizationError.USER_PERMISSION_DENIED</type> <trigger>&lt;null&gt;</trigger> <fieldPath></fieldPath> </ApiError> </reportDownloadError>`)
	client := &TestClient{
		res: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
			StatusCode: 400,
		},
	}

	auth := &Auth{
		Client: client,
	}

	rs := NewReportDownloadService(auth)
	def := ReportDefinition{}
	_, err := rs.Get(def)
	if err == nil {
		t.Fatalf("expected api error")
	}

	type ErrorCode interface {
		Code() string
	}

	expectedCode := "USER_PERMISSION_DENIED"
	if ec, ok := err.(ErrorCode); ok {
		if ec.Code() != expectedCode {
			t.Errorf("got %s, expected %s\n", ec.Code(), expectedCode)
		}
	} else {
		t.Errorf("error expected to satisfy ErrorCode interface")
	}

}

func TestReportQueryStream(t *testing.T) {
	query := `SELECT  AccountDescriptiveName, AdvertisingChannelType, Clicks, ConversionValue, Cost, Impressions, Device, ExternalCustomerId, DayOfWeek, CampaignId  FROM CAMPAIGN_PERFORMANCE_REPORT  DURING YESTERDAY`
	config := getTestConfig()
	svc := NewReportDownloadService(&config.Auth)

	report, err := svc.StreamAWQL(query, "CSV")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(report)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
}

func TestReportStreamDownloadAuthError(t *testing.T) {
	body := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?> <reportDownloadError> <ApiError> <type>AuthorizationError.USER_PERMISSION_DENIED</type> <trigger>&lt;null&gt;</trigger> <fieldPath></fieldPath> </ApiError> </reportDownloadError>`)
	client := &TestClient{
		res: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
			StatusCode: 400,
		},
	}

	auth := &Auth{
		Client: client,
	}

	rs := NewReportDownloadService(auth)
	_, err := rs.StreamAWQL("", "")
	if err == nil {
		t.Fatalf("expected api error")
	}

	type ErrorCode interface {
		Code() string
	}

	expectedCode := "USER_PERMISSION_DENIED"
	if ec, ok := err.(ErrorCode); ok {
		if ec.Code() != expectedCode {
			t.Errorf("got %s, expected %s\n", ec.Code(), expectedCode)
		}
	} else {
		t.Errorf("error expected to satisfy ErrorCode interface")
	}

}

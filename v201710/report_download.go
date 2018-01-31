package v201710

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ReportDownloadService struct {
	Auth
}

type reportDefinitionXml struct {
	*ReportDefinition
	XMLName xml.Name
}

//type is equiv to errorString eg AuthorizationError.USER_PERMISSION_DENIED
type ApiError struct {
	Type string `xml:"type"`
}

type ReportDownloadError struct {
	XMLName  xml.Name `xml:"reportDownloadError"`
	ApiError ApiError
}

func (s ApiError) Error() string {
	return s.Type
}

func (s ApiError) Code() string {
	if parts := strings.Split(s.Type, "."); len(parts) > 1 {
		return parts[1]
	}
	return s.Type
}

func NewReportDownloadService(auth *Auth) *ReportDownloadService {
	return &ReportDownloadService{Auth: *auth}
}

func (s *ReportDownloadService) Get(reportDefinition ReportDefinition) (res interface{}, err error) {
	reportDefinition.Selector.XMLName = xml.Name{baseUrl, "selector"}
	repDef := reportDefinitionXml{
		ReportDefinition: &reportDefinition,
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "reportDefinition",
		},
	}
	body, err := xml.MarshalIndent(repDef, "  ", "  ")
	if err != nil {
		return res, err
	}
	form := url.Values{}
	form.Add("__rdxml", string(body))
	resp, err := s.makeRequest(form)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		dec := xml.NewDecoder(resp.Body)
		el := &ReportDownloadError{}
		if err := dec.Decode(el); err != nil {
			return nil, err
		}
		return nil, el.ApiError
	}
	return parseReport(resp.Body)
}

func (s *ReportDownloadService) StreamAWQL(awql string, fmt string) (io.ReadCloser, error) {
	form := url.Values{}
	form.Add("__rdquery", awql)
	form.Add("__fmt", fmt)
	resp, err := s.makeRequest(form)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		dec := xml.NewDecoder(resp.Body)
		el := &ReportDownloadError{}
		if err := dec.Decode(el); err != nil {
			return nil, err
		}
		return nil, el.ApiError
	}

	return resp.Body, nil
}

func (s *ReportDownloadService) AWQL(awql string, fmt string) (interface{}, error) {
	body, err := s.StreamAWQL(awql, fmt)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseReport(body)
}

// Make our http request using the given form (re-usable for either XML or AWQL)
func (s *ReportDownloadService) makeRequest(form url.Values) (res *http.Response, err error) {
	req, err := http.NewRequest("POST", reportDownloadServiceUrl.Url, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return res, err
	}
	req.Header.Add("developerToken", s.Auth.DeveloperToken)
	req.Header.Add("clientCustomerId", s.Auth.CustomerId)
	req.Header.Add("skipReportHeader", "true")
	req.Header.Add("skipReportSummary", "true")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return s.Client.Do(req)
}

func parseReport(report io.Reader) (collection []map[string]string, err error) {
	reader := csv.NewReader(report)
	header, err := reader.Read()
	if err != nil {
		return collection, err
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return collection, err
		}
		row := make(map[string]string)
		for i := 0; i < len(record); i++ {
			column := header[i]
			row[column] = record[i]
		}
		collection = append(collection, row)
	}
	return collection, err
}

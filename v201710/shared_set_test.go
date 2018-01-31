package v201710

import (
	"bytes"
	"net/http"
	"testing"
)

type MockValidSharedSetClient struct{}
type MockReadCloser struct{}
type BufferCloser struct {
	buf *bytes.Buffer
}

func (b BufferCloser) Close() error {
	return nil
}

func (b BufferCloser) Read(p []byte) (int, error) {
	return b.buf.Read(p)
}

const GOOD_RESP = `
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><ResponseHeader xmlns="https://adwords.google.com/api/adwords/cm/v201710"><requestId>TEST_REQUEST_ID</requestId><serviceName>SharedSetService</serviceName><methodName>get</methodName><operations>1</operations><responseTime>225</responseTime></ResponseHeader></soap:Header><soap:Body><getResponse xmlns="https://adwords.google.com/api/adwords/cm/v201710"><rval><totalNumEntries>4</totalNumEntries><Page.Type>SharedSetPage</Page.Type><entries><sharedSetId>1369022798</sharedSetId><name>SQM Group 1</name><type>NEGATIVE_KEYWORDS</type><memberCount>0</memberCount><referenceCount>0</referenceCount><status>REMOVED</status></entries><entries><sharedSetId>1369022801</sharedSetId><name>SQM Group 2</name><type>NEGATIVE_KEYWORDS</type><memberCount>0</memberCount><referenceCount>0</referenceCount><status>REMOVED</status></entries><entries><sharedSetId>1369037308</sharedSetId><name>SQM Group 1</name><type>NEGATIVE_KEYWORDS</type><memberCount>160</memberCount><referenceCount>0</referenceCount><status>ENABLED</status></entries><entries><sharedSetId>1369037335</sharedSetId><name>SQM Group 2</name><type>NEGATIVE_KEYWORDS</type><memberCount>441</memberCount><referenceCount>0</referenceCount><status>ENABLED</status></entries></rval></getResponse></soap:Body></soap:Envelope>
`

func (m *MockValidSharedSetClient) Do(req *http.Request) (*http.Response, error) {
	buf := bytes.NewBuffer([]byte(GOOD_RESP))
	b := BufferCloser{buf}
	return &http.Response{
		Body: b,
	}, nil
}

func TestGetSharedSet(t *testing.T) {
	testAuth := testAuthSetup(t)
	testAuth.Client = &MockValidSharedSetClient{}

	sharedSetService := NewSharedSetService(&testAuth)
	_, count, err := sharedSetService.Get(Selector{})

	if err != nil {
		t.Fatalf("Got unexpected error: %s", err)
	}

	if count != 4 {
		t.Fatalf("Should have received 4 shared sets, got %d", count)
	}
}

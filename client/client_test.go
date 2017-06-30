package client

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"net/url"
	"net/http"
	"testing"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestClient_retryHTTP(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	endpoint, err := url.Parse("http://foo.bar")
	if err != nil {
		t.Fatal(err)
	}

	client := Client{}
	client.MaxRetries = opc.Int(5)
	// Can't use a custom transport, otherwise httpmock won't catch request
	client.httpClient = http.DefaultClient
	client.APIEndpoint = endpoint
	client.logger = opc.NewDefaultLogger()
	client.loglevel = opc.LogLevel()

	httpmock.RegisterResponder("GET", "http://foo.bar/",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(404, "mocked error message"), nil
		},
	)

	req, err := http.NewRequest("GET", "http://foo.bar/", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, reqErr := client.retryRequest(req)
	if reqErr == nil {
		t.Fatalf("Expected error, got none")
	}

	if httpmock.GetTotalCallCount() != 5 {
		t.Fatalf("Expected 5 retries, got: %d", httpmock.GetTotalCallCount())
	}
}
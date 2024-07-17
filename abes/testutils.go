package abes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// testServer wraps all logic used to simulate HTTP responses, including
// network failure.
type testServer struct {
	server         *httptest.Server
	client         *http.Client
	roundTripper   *customRoundTripper
	networkFailure bool
}

// newTestServer creates and configures the server and the associated client.
func newTestServer() *testServer {
	ts := new(testServer)
	ts.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ts.networkFailure {
			http.Error(w, "network failure", http.StatusServiceUnavailable)
			return
		}
		fmt.Println(r.URL)
		fmt.Fprintln(w, "Hello, client")
	}))
	ts.roundTripper = &customRoundTripper{
		originalTransport: http.DefaultTransport,
		testServerURL:     ts.server.URL,
	}
	ts.client = &http.Client{
		Transport: ts.roundTripper,
	}
	return ts
}

func (ts *testServer) close() {
	ts.server.Close()
}

// simulateNetworkFailure enables or disables a network failure simulation.
func (ts *testServer) simulateNetworkFailure(enable bool) {
	ts.networkFailure = enable
	ts.roundTripper.networkFailure = enable
}

// CustomRoundTripper redirects requests toward the test server.
type customRoundTripper struct {
	originalTransport http.RoundTripper
	testServerURL     string
	networkFailure    bool
}

func (c *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// if c.networkFailure {
	// 	return nil, errors.New("network failure")
	// }

	// redirection
	req.URL.Host = c.testServerURL[7:]
	req.URL.Scheme = "http"
	req.Host = c.testServerURL[7:]
	return c.originalTransport.RoundTrip(req)
}

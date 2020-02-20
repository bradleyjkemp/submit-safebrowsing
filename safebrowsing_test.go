package safebrowsing

import (
	"io/ioutil"
	"net/http"
	"testing"
)

type mockRoundTripper struct{}

func (mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(nil),
	}, nil
}

func TestSubmit(t *testing.T) {
	sub := Submitter{&http.Client{
		Transport: mockRoundTripper{},
	}}
	err := sub.Submit(Report{
		URL:        "tesco-my-accounts.000webhostapp.com",
		Screenshot: nil,
		//DOM:        "",
		Flags: []Flag{NotTopSite, NotVisitedBefore},
	})
	if err != nil {
		t.Fatal(err)
	}
}

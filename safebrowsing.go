package safebrowsing

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// Based on ChromeExtensionClientRequest in:
// https://github.com/chromium/suspicious-site-reporter/blob/310bd04ffa5fa750100be6c6096473dd824857b6/extension/client_request.proto#L21
type Report struct {
	// The suspicious URL
	URL string // JSPB field 1

	// A PNG screenshot
	Screenshot []byte // JSPB field 3

	// The DOM of the suspicious site
	DOM string // JSPB field 4

	// JSPB field 5 ReferrerChain not yet supported

	// A set of flags as to why this is suspicious
	Flags []Flag // JSPB field 6
}

type Flag string

const (
	IsIDN                                Flag = "isIDN"                                // Domain uses uncommon characters
	LongSubdomains                       Flag = "longSubdomains"                       // Unusually long subdomains
	NotTopSite                           Flag = "notTopSite"                           // Site not in top 5k sites
	NotVisitedBefore                     Flag = "notVisitedBefore"                     // Haven't visited site in the last 3 months
	ManySubdomains                       Flag = "manySubdomains"                       // Unusually many subdomains
	RedirectsThroughSuspiciousTld        Flag = "redirectsThroughSuspiciousTld"        // Site redirected through a TLD potentially associated with abuse
	RedirectsFromOutsideProgramOrWebmail Flag = "redirectsFromOutsideProgramOrWebmail" // Visit maybe initiated from outside program or webmail
	UrlShortenerRedirects                Flag = "urlShortenerRedirects"                // Has multiple redirects through URL shorteners
)

// A super hacky jspb marshaller for this request
func (r Report) MarshalJSON() ([]byte, error) {
	jspbReport := make([]interface{}, 6)
	jspbReport[0] = r.URL
	jspbReport[2] = base64.StdEncoding.EncodeToString(r.Screenshot)
	jspbReport[3] = r.DOM
	jspbReport[5] = r.Flags
	return json.Marshal(jspbReport)
}

func Submit(report Report) error {
	return Submitter{HTTPClient: http.DefaultClient}.Submit(report)
}

type Submitter struct {
	HTTPClient *http.Client
}

const safeBrowsingSubmitAPI = "https://safebrowsing.google.com/safebrowsing/clientreport/crx-report"

func (s Submitter) Submit(report Report) error {
	marshalled, err := json.Marshal(report)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, safeBrowsingSubmitAPI, bytes.NewReader(marshalled))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("User-Agent", "github.com/bradleyjkemp/safebrowsing")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	spew.Dump(resp)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}
	return nil
}

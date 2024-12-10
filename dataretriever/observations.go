package dataretriever

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func FetchLatestObservationsXMLByProvince(province string) (io.ReadCloser, error) {
	now := time.Now().UTC()
	timestamp := now.Format("2006010215")
	url := fmt.Sprintf("https://dd.weather.gc.ca/observations/xml/%s/hourly/hourly_%s_%s_e.xml", province, strings.ToLower(province), timestamp)
	fmt.Printf("Fetching URL: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}
	return resp.Body, nil
}

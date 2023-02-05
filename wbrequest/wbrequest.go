package wbrequest

import (
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	kRetryIntervalSec = 3
	kAttemptsCount = 5
)

func SendWithRetries(method string, url string, headers map[string]string) (body []byte, status_code int, err error) {
	log.Debugf("http request method: %v url: %v", method, url)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Errorf("http.NewRequest: %v\n", err)
		return nil, 0, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-GB,en;q=0.9")
	req.Header.Add("Cache-Control", "no-store")

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	r, err := DoWithRetries(req)
	if err != nil {
		log.Errorf("Could not send request: %s", err)
		return nil, 0, err
	}
	defer r.Body.Close()

	body, err = io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Could not read response body: %s", err)
		return nil, 0, err
	}
	return body, r.StatusCode, err
}

func DoWithRetries(request *http.Request) (response *http.Response, err error){
	client := &http.Client{}
	for attempt := 0; attempt < kAttemptsCount; attempt++ {
		response, err = client.Do(request)
		if err != nil {
			return nil, err
		}
		log.Debugf("status code: %v\n", response.StatusCode)
		if response.StatusCode == http.StatusOK {
			break
		}
		log.Debugf("retry after %d sec\n", kRetryIntervalSec)
		time.Sleep(time.Duration(kRetryIntervalSec) * time.Second)
	}
	return response, nil
}
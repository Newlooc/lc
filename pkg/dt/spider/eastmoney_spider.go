package spider

import (
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	defaultHeader = map[string][]string{
		"User-Agent": {
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36",
		},
	}
)

type Visit struct {
	httpClient *http.Client
	Header     http.Header
	l          *sync.Mutex
}

func NewVisit() *Visit {
	return &Visit{
		httpClient: &http.Client{},
		Header:     defaultHeader,
		l:          &sync.Mutex{},
	}
}

func (v *Visit) Do(url string, dryrun bool) (http.Header, []byte, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		log.WithError(err).Errorf("Failed to parse %s", url)
		return nil, nil, err
	}

	req := &http.Request{}
	req.Header = v.Header
	req.URL = parsedURL

	if dryrun {
		log.Infof("Run URL %s", url)
		return nil, nil, nil
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Errorf("Failed to visit %s", url)
		return nil, nil, err
	}

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Errorf("Failed to read IO %s", url)
		return nil, nil, err
	}

	return resp.Header, ret, nil
}

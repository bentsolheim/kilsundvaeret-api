package app

import (
	"encoding/json"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func ConfigureLogging(logLevel string) error {
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:    true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return stacktrace.Propagate(err, "error while parsing log level %s", logLevel)
	}
	log.SetLevel(level)
	return nil
}

func HttpGetJson(url string, response interface{}) error {
	body, err := HttpGet(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return stacktrace.Propagate(err, "error while unmarshalling response")
	}
	return nil
}

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed getting data from [%s]", url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error reading response body")
	}
	return body, nil
}

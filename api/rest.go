package api

import (
	"github.com/0x111/sn-edit/conf"
	log "github.com/sirupsen/logrus"
)

func Get(url string) ([]byte, error) {
	restClient := conf.GetClient()

	resp, err := restClient.R().
		Get(url)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was a problem getting the entry from the instance! Please try again later!")
		return nil, err
	}

	if resp.StatusCode() != 200 {
		log.WithFields(log.Fields{"status_code": resp.StatusCode()}).Error("We received a HTTP Error Code from the Instance. Please check your config file and try again.")
		return nil, err
	}

	return resp.Body(), nil
}

func Put(url string, body interface{}) ([]byte, error) {
	restClient := conf.GetClient()

	resp, err := restClient.R().
		SetBody(body).
		Put(url)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was a problem uploading the entry from the instance! Please try again later!")
		return nil, err
	}

	return resp.Body(), nil
}

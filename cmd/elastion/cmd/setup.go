package cmd

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ElasticClient *elasticsearch7.Client
var log = logrus.New()

func setup() {

	if prettyLog {
		log.SetFormatter(&logrus.TextFormatter{})
	} else {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	log.SetOutput(os.Stdout)

	esCfg := elasticsearch7.Config{
		Addresses: viper.GetStringSlice("esAddresses"),
		Username:  viper.GetString("esUsername"),
		Password:  viper.GetString("esPassword"),
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   viper.GetInt("esMaxIdleConns"),
			ResponseHeaderTimeout: viper.GetDuration("esDialTimeout"),
			DialContext:           (&net.Dialer{Timeout: viper.GetDuration("esResponseHeaderTimeout")}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
	}

	var err error = nil

	ElasticClient, err = elasticsearch7.NewClient(esCfg)

	if err != nil {
		log.Fatalf("Error when starting elasticsearch: %v", err)
	}

}

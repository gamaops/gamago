package elastic

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Options struct {
	Addresses             []string
	Username              string
	Password              string
	MaxIdleConnsPerHost   int
	ResponseHeaderTimeout time.Duration
	DialContext           time.Duration
}

type Elastic struct {
	Options *Options
	Client  *elasticsearch7.Client
	Logger  *logrus.Logger
}

func StartElasticClient(elastic *Elastic) error {

	esCfg := elasticsearch7.Config{
		Addresses: elastic.Options.Addresses,
		Username:  elastic.Options.Username,
		Password:  elastic.Options.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   elastic.Options.MaxIdleConnsPerHost,
			ResponseHeaderTimeout: elastic.Options.ResponseHeaderTimeout,
			DialContext:           (&net.Dialer{Timeout: elastic.Options.DialContext}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
	}

	client, err := elasticsearch7.NewClient(esCfg)

	if err != nil {
		return err
	}

	elastic.Client = client

	return nil

}

func SetupViper() {
	viper.SetDefault("esAddresses", "http://localhost:9200")
	viper.BindEnv("esAddresses", "ES_ADDRESSES")
	viper.SetDefault("esUsername", "")
	viper.BindEnv("esUsername", "ES_USERNAME")
	viper.SetDefault("esPassword", "")
	viper.BindEnv("esPassword", "ES_PASSWORD")
	viper.SetDefault("esMaxIdleConns", "10")
	viper.BindEnv("esMaxIdleConns", "ES_MAX_IDLE_CONNS")
	viper.SetDefault("esDialTimeout", "5s")
	viper.BindEnv("esDialTimeout", "ES_DIAL_TIMEOUT")
	viper.SetDefault("esResponseHeaderTimeout", "5s")
	viper.BindEnv("esResponseHeaderTimeout", "ES_RESPONSE_HEADER_TIMEOUT")
}

type ErrElasticResponse struct {
	Raw      map[string]interface{}
	Response *esapi.Response
}

func (e *ErrElasticResponse) Error() string {
	return fmt.Sprintf("error on Elasticsearch response: %v", e.Raw)
}

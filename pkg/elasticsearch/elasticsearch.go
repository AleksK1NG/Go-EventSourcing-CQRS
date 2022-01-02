package elasticsearch

import (
	"context"
	v7 "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

type Config struct {
	URL         string `mapstructure:"url"`
	Sniff       bool   `mapstructure:"sniff"`
	Gzip        bool   `mapstructure:"gzip"`
	Explain     bool   `mapstructure:"explain"`
	FetchSource bool   `mapstructure:"fetchSource"`
	Version     bool   `mapstructure:"version"`
	Pretty      bool   `mapstructure:"pretty"`
}

func NewElasticClient(ctx context.Context, cfg Config) (*v7.Client, error) {
	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := v7.NewClient(
		v7.SetURL(cfg.URL),
		v7.SetSniff(cfg.Sniff),
		v7.SetGzip(cfg.Gzip),
	)
	if err != nil {
		return nil, errors.Wrap(err, "v7.NewClient")
	}

	//// Ping the Elasticsearch server to get e.g. the version number
	//info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	//if err != nil {
	//	return nil, errors.Wrap(err, "client.Ping")
	//}
	//fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	//
	//// Getting the ES version number is quite common, so there's a shortcut
	//esVersion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	//if err != nil {
	//	return nil, errors.Wrap(err, "client.ElasticsearchVersion")
	//}
	//fmt.Printf("Elasticsearch version %s\n", esVersion)

	return client, nil
}

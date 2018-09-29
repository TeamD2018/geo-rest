package services

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"log"
	"os"
	"testing"
)

var client *elastic.Client

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("bitnami/elasticsearch", "6.4.1", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		addr := fmt.Sprintf("http://localhost:%s", resource.GetPort("9200/tcp"))

		var err error
		client, err = elastic.NewClient(elastic.SetURL(addr), elastic.SetSniff(false))
		if err != nil {
			return err
		}

		_, _, err = client.Ping(addr).Do(context.Background())

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

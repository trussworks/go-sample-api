package integration

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/facebookgo/clock"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/trussworks/httpbaselinetest"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"bin/bork/pkg/appconfig"
	"bin/bork/pkg/server/httpserver"
	"bin/bork/pkg/sources/postgres"
)

func TestBaselines(t *testing.T) {
	config := viper.New()
	config.AutomaticEnv()

	origAppConfig, err := appconfig.NewAppConfig(config)
	if err != nil {
		t.Fatalf("Error initializing config: %s", err)
	}

	dataSourceName := postgres.BuildDataSourceName(&origAppConfig)

	txdb.Register("pgx", "postgres", dataSourceName)

	testAppConfig := origAppConfig
	testAppConfig.DBDriver = "pgx"
	testAppConfig.Clock = clock.NewMock()
	setupFunc := func(name string, btest *httpbaselinetest.HttpBaselineTest) error {
		// seed the source based on the name of the test
		h := fnv.New64a()
		h.Write([]byte(name))
		uuid.SetRand(rand.New(rand.NewSource(int64(h.Sum64()))))

		testAppConfig.DBName = "txdb_" + httpbaselinetest.NormalizeTestName(name)

		observer, logs := observer.New(zap.InfoLevel)
		testAppConfig.Logger = zap.New(observer)

		server, err := httpserver.NewServer(&testAppConfig)
		if err != nil {
			return fmt.Errorf("Error initializing server: %s", err)
		}
		btest.Handler = server
		btest.Db = server.Db()
		btest.Custom = logs
		return nil
	}

	teardownFunc := func(t *testing.T, btest *httpbaselinetest.HttpBaselineTest) error {
		if t.Failed() && btest.Custom != nil {
			logs := btest.Custom.(*observer.ObservedLogs)
			for _, log := range logs.All() {
				t.Log(log)
			}
		}
		return nil
	}

	bts := httpbaselinetest.NewDefaultHttpBaselineTestSuite(t)

	bts.Run("GET v1 dog without auth", httpbaselinetest.HttpBaselineTest{
		Setup:    setupFunc,
		Teardown: teardownFunc,
		Method:   http.MethodGet,
		Path:     "/api/v1/dog/00000000-0000-0000-0000-000000000000",
		Seed:     "chihuahua.seed.yml",
	})

	bts.Run("GET v1 dog with auth", httpbaselinetest.HttpBaselineTest{
		Setup:    setupFunc,
		Teardown: teardownFunc,
		Method:   http.MethodGet,
		Path:     "/api/v1/dog/00000000-0000-0000-0000-000000000000",
		Headers: map[string]string{
			"Authorization": "Owner",
		},
		Seed: "chihuahua.seed.yml",
	})

	bts.Run("GET v1 dog missing", httpbaselinetest.HttpBaselineTest{
		Setup:    setupFunc,
		Teardown: teardownFunc,
		Method:   http.MethodGet,
		Path:     "/api/v1/dog/00000000-0000-0000-0000-000000000000",
		Headers: map[string]string{
			"Authorization": "Owner",
		},
	})

	bts.Run("POST v1 dog with auth", httpbaselinetest.HttpBaselineTest{
		Setup:    setupFunc,
		Teardown: teardownFunc,
		Method:   http.MethodPost,
		Path:     "/api/v1/dog",
		Body: map[string]string{
			"name":      "Lola",
			"breed":     "Chihuahua",
			"birthDate": testAppConfig.Clock.Now().Format(time.RFC3339),
		},
		Headers: map[string]string{
			"Authorization": "Owner",
			"Content-Type":  "application/json",
		},
		Tables: []string{"dog"},
	})

	emptyMap := make(map[string]string, 0)
	bts.Run("GraphQL fetch all dogs", httpbaselinetest.HttpBaselineTest{
		Setup:    setupFunc,
		Teardown: teardownFunc,
		Method:   http.MethodPost,
		Path:     "/graphql/query",
		Body: map[string]interface{}{
			"operationName": nil,
			"variables":     emptyMap,
			"query": `
{
  dogs {
    id, name, breed, birthDate, owner { id }
  }
}`,
		},
		Headers: map[string]string{
			"Authorization": "Owner",
			"Content-Type":  "application/json",
		},
		Seed: "chihuahua.seed.yml",
	})

	bts.Run("GraphQL fetch single dog", httpbaselinetest.HttpBaselineTest{
		Setup:    setupFunc,
		Teardown: teardownFunc,
		Method:   http.MethodPost,
		Path:     "/graphql/query",
		Body: map[string]interface{}{
			"operationName": nil,
			"variables":     emptyMap,
			"query": `
{
  dog(dogId: "00000000-0000-0000-0000-000000000000") {
    id, name, breed, birthDate, owner { id }
  }
}`,
		},
		Headers: map[string]string{
			"Authorization": "Owner",
			"Content-Type":  "application/json",
		},
		Seed: "chihuahua.seed.yml",
	})

	bts.Run("GraphQL create dog", httpbaselinetest.HttpBaselineTest{
		Setup:    setupFunc,
		Teardown: teardownFunc,
		Method:   http.MethodPost,
		Path:     "/graphql/query",
		Body: map[string]interface{}{
			"operationName": nil,
			"variables":     emptyMap,
			"query": `
mutation {
  createDog(input: {
    birthDate: "` + testAppConfig.Clock.Now().Format(time.RFC3339) + `",
    name: "Lola",
    breed: CHIHUAHUA
  }) {
    id, name, breed, birthDate, owner { id }
  }
}`,
		},
		Headers: map[string]string{
			"Authorization": "Owner",
			"Content-Type":  "application/json",
		},
		Tables: []string{"dog"},
	})
}

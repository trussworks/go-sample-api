package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/google/uuid"

	"bin/bork/pkg/models"
)

func (s IntegrationTestSuite) TestDogEndpoints() {
	apiURL, err := url.Parse(s.server.URL)
	s.NoError(err, "failed to parse URL")
	apiURL.Path = path.Join(apiURL.Path, "/api/v1")
	dogURL, err := url.Parse(apiURL.String())
	s.NoError(err, "failed to parse URL")
	dogURL.Path = path.Join(dogURL.Path, "/dog")

	client := &http.Client{}

	s.Run("POST will fail with no Authorization", func() {
		req, err := http.NewRequest(http.MethodPost, dogURL.String(), bytes.NewBufferString(""))
		s.NoError(err)
		resp, err := client.Do(req)

		s.NoError(err)
		s.Equal(http.StatusUnauthorized, resp.StatusCode)
	})

	postDog := models.Dog{}
	owner := "Owner"
	s.Run("POST will succeed", func() {
		body, err := json.Marshal(map[string]string{
			"name":      "Lola",
			"breed":     "Chihuahua",
			"birthDate": s.clock.Now().Format(time.RFC3339),
		})
		s.NoError(err)
		req, err := http.NewRequest(http.MethodPost, dogURL.String(), bytes.NewBuffer(body))
		s.NoError(err)
		req.Header.Set("Authorization", owner)

		resp, err := client.Do(req)

		s.NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)
		actualBody, err := ioutil.ReadAll(resp.Body)
		s.NoError(err)
		err = json.Unmarshal(actualBody, &postDog)
		s.NoError(err)
		s.NotZero(postDog.ID)
	})

	getURL, err := url.Parse(dogURL.String())
	s.NoError(err, "failed to parse URL")
	getURL.Path = path.Join(getURL.Path, postDog.ID.String())

	s.Run("GET will fetch the dog just saved", func() {
		req, err := http.NewRequest(http.MethodGet, getURL.String(), nil)
		s.NoError(err)
		req.Header.Set("Authorization", postDog.OwnerID)

		resp, err := client.Do(req)

		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(http.StatusOK, resp.StatusCode)
		actualBody, err := ioutil.ReadAll(resp.Body)
		s.NoError(err)
		getDog := models.Dog{}
		err = json.Unmarshal(actualBody, &getDog)
		s.NoError(err)
		s.Equal(postDog.ID, getDog.ID)
		s.Equal(postDog.Name, getDog.Name)
		s.Equal(postDog.Breed, getDog.Breed)
		s.Equal(postDog.OwnerID, getDog.OwnerID)
		s.True(postDog.BirthDate.Equal(getDog.BirthDate))
	})

	s.Run("GET will fail with no Authorization", func() {
		req, err := http.NewRequest(http.MethodGet, getURL.String(), bytes.NewBufferString(""))
		s.NoError(err)
		resp, err := client.Do(req)

		s.NoError(err)
		s.Equal(http.StatusUnauthorized, resp.StatusCode)
	})

	s.Run("GET will fail with wrong Owner", func() {
		req, err := http.NewRequest(http.MethodGet, getURL.String(), bytes.NewBufferString(""))
		s.NoError(err)
		req.Header.Set("Authorization", "Other Owner")
		resp, err := client.Do(req)

		s.NoError(err)
		s.Equal(http.StatusUnauthorized, resp.StatusCode)
	})

	s.Run("PUT will fail with no Authorization", func() {
		req, err := http.NewRequest(http.MethodPut, dogURL.String(), bytes.NewBufferString(""))
		s.NoError(err)
		resp, err := client.Do(req)

		s.NoError(err)
		s.Equal(http.StatusUnauthorized, resp.StatusCode)
	})

	s.Run("PUT will succeed", func() {
		postDog.Name = "Lolita"
		body, err := json.Marshal(postDog)
		s.NoError(err)
		req, err := http.NewRequest(http.MethodPut, dogURL.String(), bytes.NewBuffer(body))
		s.NoError(err)
		req.Header.Set("Authorization", owner)

		resp, err := client.Do(req)

		s.NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)
		actualBody, err := ioutil.ReadAll(resp.Body)
		s.NoError(err)
		putDog := models.Dog{}
		err = json.Unmarshal(actualBody, &putDog)
		s.NoError(err)
		s.Equal("Lolita", putDog.Name)
	})

	s.Run("PUT will fail on unknown dog", func() {
		body, err := json.Marshal(map[string]string{
			"id":        uuid.New().String(),
			"name":      "Lola",
			"breed":     "Chihuahua",
			"birthDate": s.clock.Now().Format(time.RFC3339),
		})
		s.NoError(err)
		req, err := http.NewRequest(http.MethodPut, dogURL.String(), bytes.NewBuffer(body))
		s.NoError(err)
		req.Header.Set("Authorization", owner)

		resp, err := client.Do(req)

		s.NoError(err)
		s.Equal(http.StatusNotFound, resp.StatusCode)
	})

	s.Run("GET all will fetch the dog saved", func() {
		dogsURL, err := url.Parse(apiURL.String())
		s.NoError(err, "failed to parse URL")
		dogsURL.Path = path.Join(dogsURL.Path, "/dogs")

		req, err := http.NewRequest(http.MethodGet, dogsURL.String(), nil)
		s.NoError(err)
		req.Header.Set("Authorization", postDog.OwnerID)

		resp, err := client.Do(req)

		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(http.StatusOK, resp.StatusCode)
		actualBody, err := ioutil.ReadAll(resp.Body)
		s.NoError(err)
		getDogs := models.Dogs{}
		err = json.Unmarshal(actualBody, &getDogs)
		s.NoError(err)
		postDogCount := 0
		for _, d := range getDogs {
			if d.ID == postDog.ID {
				postDogCount++
			}
		}
		s.Equal(postDogCount, 1)
	})
}

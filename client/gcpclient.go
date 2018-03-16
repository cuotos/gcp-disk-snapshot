package client

import (
	"context"
	"errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"io/ioutil"
	"net/http"
)

func NewComputeClient() (*compute.Service, error) {

	requiredScopes := []string{compute.ComputeScope}

	client, err := google.DefaultClient(context.Background(), requiredScopes...)
	if err != nil {
		return nil, err
	}

	return compute.New(client)
}

func EstablishProjectId() (string, error) {
	//curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/project/project-id
	projectIdMeta := `http://metadata.google.internal/computeMetadata/v1/project/project-id`

	req, err := http.NewRequest("GET", projectIdMeta, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Metadata-Flavor", "Google")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}
	res.Body.Close()

	contentString := string(body)

	if contentString == "" {
		return "", errors.New("unable to establish projectid")
	}

	return contentString, nil
}

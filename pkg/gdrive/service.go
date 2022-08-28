package gdrive

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/trichner/oauthflows"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	clientSecretFile = "client_secret.json"
	credentialsFile  = "credentials.json"
)

var scopes = []string{
	"https://www.googleapis.com/auth/drive",
	"https://www.googleapis.com/auth/drive.file",
	"https://www.googleapis.com/auth/drive.readonly",
}

type DriveService struct {
	service *drive.Service
}

type Drive struct {
	// Id: The ID of this Drive which is also the ID of the top level
	// folder of this Team Drive.
	Id string `json:"id,omitempty"`

	// Kind: Identifies what kind of resource this is.
	Kind string `json:"kind,omitempty"`

	// Name: The name of this Drive.
	Name string `json:"name,omitempty"`
}

func NewDriveService() (*DriveService, error) {
	var options []option.ClientOption

	log.Printf("reading %q", clientSecretFile)
	client, err := newOAuthClientFormClientSecret()
	if err == nil {
		log.Printf("creating api client from %q", clientSecretFile)
		options = append(options, option.WithHTTPClient(client))
	} else {
		log.Printf("reading %q", credentialsFile)
		j, err := readCredentialsJson()
		if err != nil {
			return nil, err
		}

		log.Printf("creating api client from %q", credentialsFile)
		options = append(options, option.WithCredentialsJSON(j))
	}

	service, err := drive.NewService(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	return &DriveService{service: service}, nil
}

func newOAuthClientFormClientSecret() (*http.Client, error) {
	return oauthflows.NewClient(oauthflows.WithClientSecretsFile(clientSecretFile, scopes))
}

func readCredentialsJson() ([]byte, error) {
	return ioutil.ReadFile(credentialsFile)
}

func (s *DriveService) ListDrives() ([]*Drive, error) {
	resp, err := s.service.Drives.List().PageSize(100).Do()
	if err != nil {
		return nil, fmt.Errorf("cannot list drives: %w", err)
	}

	drives := make([]*Drive, len(resp.Drives))
	for i, d := range resp.Drives {
		drives[i] = &Drive{
			Id:   d.Id,
			Kind: d.Kind,
			Name: d.Name,
		}
	}
	return drives, nil
}

package gdrive

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/trichner/oauthflows"
	"github.com/trichner/toolbox/cmd/tb/cfg"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var scopes = []string{
	"https://www.googleapis.com/auth/drive",
	"https://www.googleapis.com/auth/drive.file",
	"https://www.googleapis.com/auth/drive.readonly",
}

const clientCredentialsFile = "client_credentials.json"

func Exec(ctx context.Context, args []string) {
	config := cfg.FromContext(ctx)
	httpClient, err := newHttpClient(config, scopes)
	if err != nil {
		panic(err)
	}
	service, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		panic(err)
	}

	file := "test"
	fd, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	lastUpdated := time.Now()
	_, err = service.Files.Create(&drive.File{
		Name:    file,
		Parents: nil,
	}).Media(fd).ProgressUpdater(func(current, total int64) {
		if time.Since(lastUpdated) > 2*time.Second {
			lastUpdated = time.Now()
		}
	}).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
}

func newHttpClient(cfg cfg.ConfigProvider, scopes []string) (*http.Client, error) {
	slurp, err := cfg.ReadFile(clientCredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", clientCredentialsFile, err)
	}

	config, err := google.ConfigFromJSON(slurp, scopes...)
	if err != nil {
		return nil, err
	}

	return oauthflows.NewClient(oauthflows.WithConfig(config))
}

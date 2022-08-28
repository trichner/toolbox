package sheets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trichner/oauthflows"
	"golang.org/x/oauth2/google"
	. "google.golang.org/api/option"
	googlesheets "google.golang.org/api/sheets/v4"
)

const clientSecretFile = "client_secret.json"

var scopes = []string{
	"https://www.googleapis.com/auth/drive",
	"https://www.googleapis.com/auth/drive.file",
	"https://www.googleapis.com/auth/drive.readonly",
	"https://www.googleapis.com/auth/spreadsheets",
	"https://www.googleapis.com/auth/spreadsheets.readonly",
}

type SheetsService interface {
	CreateSpreadSheet(title string) (SpreadsheetOps, error)
	GetSpreadSheet(id string) (SpreadsheetOps, error)
}

type sheetsService struct {
	service *googlesheets.Service
}

type SpreadSheet struct {
	Id     string
	Sheets []*Sheet
}

type Sheet struct {
	Id    int64
	Title string
	Index int64
}

func NewSheetService(ctx context.Context) (SheetsService, error) {
	var err error
	var client *http.Client

	if client, err = google.DefaultClient(ctx, scopes...); err == nil {
		// client already set
	} else if client, err = oauthflows.NewClient(oauthflows.WithClientSecretsFile(clientSecretFile, scopes), oauthflows.WithFileTokenStore()); err == nil {
		// client already set
	} else {
		return nil, fmt.Errorf("cannot initialize oauth client: %w", err)
	}

	service, err := googlesheets.NewService(ctx, WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("cannot create service: %w", err)
	}

	return &sheetsService{service: service}, nil
}

func (s *sheetsService) GetSpreadSheet(id string) (SpreadsheetOps, error) {
	res, err := s.service.Spreadsheets.Get(id).Do()
	if err != nil {
		return nil, err
	}
	return &spreadsheetOps{
		service:     s.service,
		spreadsheet: res,
	}, nil
}

func (s *sheetsService) CreateSpreadSheet(title string) (SpreadsheetOps, error) {
	ss := &googlesheets.Spreadsheet{
		Properties: &googlesheets.SpreadsheetProperties{Title: title},
	}
	res, err := s.service.Spreadsheets.Create(ss).Do()
	if err != nil {
		return nil, fmt.Errorf("cannot create spreadsheet: %w", err)
	}
	return &spreadsheetOps{
		service:     s.service,
		spreadsheet: res,
	}, nil
}
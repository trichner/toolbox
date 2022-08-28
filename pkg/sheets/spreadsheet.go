package sheets

import (
	"errors"
	"fmt"

	googlesheets "google.golang.org/api/sheets/v4"
)

var ErrNotFound = errors.New("not found")

type CreateSheetOptions struct {
	Title string
}

type SpreadsheetOps interface {
	CreateSheet(opts *CreateSheetOptions) (SheetOps, error)
	FirstSheet() (SheetOps, error)
	SheetByIndex(index int64) (SheetOps, error)
	SheetById(id int64) (SheetOps, error)
	SheetByTitle(name string) (SheetOps, error)
	Get() (*SpreadSheet, error)
}

type spreadsheetOps struct {
	service     *googlesheets.Service
	spreadsheet *googlesheets.Spreadsheet
}

func (s *spreadsheetOps) CreateSheet(opts *CreateSheetOptions) (SheetOps, error) {
	req := &googlesheets.AddSheetRequest{
		Properties: &googlesheets.SheetProperties{
			Hidden: false,
			Index:  0,
			Title:  opts.Title,
		},
	}

	breq := &googlesheets.BatchUpdateSpreadsheetRequest{Requests: []*googlesheets.Request{{AddSheet: req}}}

	res, err := s.service.Spreadsheets.BatchUpdate(s.spreadsheet.SpreadsheetId, breq).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to add sheet %q to %q: %w", s.spreadsheet.SpreadsheetId, opts.Title, err)
	}

	props := res.Replies[0].AddSheet.Properties

	return s.toSheetOps(props), nil
}

func (s *spreadsheetOps) FirstSheet() (SheetOps, error) {
	return s.SheetByIndex(0)
}

func (s *spreadsheetOps) SheetByIndex(index int64) (SheetOps, error) {
	return s.toSheetOpsWithErr(s.filteredSheets(func(sheet *googlesheets.SheetProperties) bool {
		return sheet.Index == index
	}))
}

func (s *spreadsheetOps) SheetById(id int64) (SheetOps, error) {
	return s.toSheetOpsWithErr(s.filteredSheets(func(sheet *googlesheets.SheetProperties) bool {
		return sheet.SheetId == id
	}))
}

func (s *spreadsheetOps) SheetByTitle(title string) (SheetOps, error) {
	return s.toSheetOpsWithErr(s.filteredSheets(func(sheet *googlesheets.SheetProperties) bool {
		return sheet.Title == title
	}))
}

func (s *spreadsheetOps) Get() (*SpreadSheet, error) {
	err := s.refresh()
	if err != nil {
		return nil, err
	}

	return &SpreadSheet{Id: s.spreadsheet.SpreadsheetId, Sheets: mapSheets(s.spreadsheet.Sheets)}, nil
}

func (s *spreadsheetOps) filteredSheets(predicate func(p *googlesheets.SheetProperties) bool) (*googlesheets.SheetProperties, error) {
	sheets, err := s.getSheets()
	if err != nil {
		return nil, err
	}

	for _, sheet := range sheets {
		props := sheet.Properties
		if predicate(props) {
			return props, nil
		}
	}
	return nil, ErrNotFound
}

func (s *spreadsheetOps) spreadsheetId() string {
	return s.spreadsheet.SpreadsheetId
}

func (s *spreadsheetOps) getSheets() ([]*googlesheets.Sheet, error) {
	err := s.refresh()
	if err != nil {
		return nil, err
	}
	return s.spreadsheet.Sheets, nil
}

func (s *spreadsheetOps) refresh() error {
	res, err := s.service.Spreadsheets.Get(s.spreadsheet.SpreadsheetId).Do()
	if err != nil {
		return fmt.Errorf("cannot create spreadsheet: %w", err)
	}

	s.spreadsheet = res

	return nil
}

func (s *spreadsheetOps) toSheetOpsWithErr(sheet *googlesheets.SheetProperties, err error) (*sheetOps, error) {
	return s.toSheetOps(sheet), err
}

func (s *spreadsheetOps) toSheetOps(sheet *googlesheets.SheetProperties) *sheetOps {
	if sheet == nil {
		return nil
	}
	return &sheetOps{
		spreadsheetOps: s,
		sheetId:        sheet.SheetId,
	}
}

func mapSheets(from []*googlesheets.Sheet) []*Sheet {
	sts := make([]*Sheet, 0, len(from))
	for _, sheet := range from {
		sts = append(sts, toSheet(sheet.Properties))
	}
	return sts
}

func toSheet(sheet *googlesheets.SheetProperties) *Sheet {
	if sheet == nil {
		return nil
	}

	return &Sheet{
		Id:    sheet.SheetId,
		Title: sheet.Title,
		Index: sheet.Index,
	}
}

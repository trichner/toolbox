package json2sheet

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/trichner/toolbox/pkg/sheets"
)

type SheetWriter interface {
	UpdateValues(data [][]string) error
}

func WriteToNewSheet(r io.Reader) (*url.URL, error) {
	ctx := context.Background()

	svc, err := sheets.NewSheetService(ctx)
	if err != nil {
		return nil, err
	}

	ss, err := svc.CreateSpreadSheet("json2sheet")
	if err != nil {
		return nil, err
	}

	sheet, err := ss.FirstSheet()
	if err != nil {
		return nil, err
	}

	err = WriteObjectsTo(sheet, r)
	if err != nil {
		return nil, err
	}

	info, err := ss.Get()
	if err != nil {
		return nil, err
	}

	raw := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit#gid=0", info.Id)
	return url.Parse(raw)
}

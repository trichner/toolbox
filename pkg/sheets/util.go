package sheets

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

// ParseSpreadsheetUrl parses a URL to a spreadsheet such as: https://docs.google.com/spreadsheets/d/1dAN8MO9NDVPqVIoOxC9H_j4Ir5c1viQ97igGdXOyXsU/edit#gid=886605725
func ParseSpreadsheetUrl(u string) (string, int64, error) {

	parsed, err := url.Parse(u)
	if err != nil {
		return "", -1, err
	}

	const googleDocsHost = "docs.google.com"
	if parsed.Host != googleDocsHost {
		return "", -1, fmt.Errorf("unexpected host '%s', expected '%s'", parsed.Host, googleDocsHost)
	}

	const httpsScheme = "https"
	if parsed.Scheme != "https" {
		return "", -1, fmt.Errorf("unexpected scheme '%s', expected '%s'", parsed.Scheme, httpsScheme)
	}

	pathPattern := regexp.MustCompile("^/spreadsheets/d/([-_A-Za-z0-9]+)/edit$")
	matches := pathPattern.FindStringSubmatch(parsed.Path)
	if matches == nil {
		return "", -1, fmt.Errorf("can't find spreadsheetId in path: '%s'", parsed.Path)
	}
	spreadsheetId := matches[1]

	q, err := url.ParseQuery(parsed.Fragment)
	if err != nil {
		return "", -1, fmt.Errorf("can't parse fragment '%s': %w", parsed.Fragment, err)
	}

	const queryParamGid = "gid"
	rawSheetId := q.Get(queryParamGid)
	if rawSheetId == "" {
		return "", -1, fmt.Errorf("can't find '%s' in '%s'", queryParamGid, parsed.Fragment)
	}

	sheetId, err := strconv.ParseInt(rawSheetId, 10, 64)

	return spreadsheetId, sheetId, err
}

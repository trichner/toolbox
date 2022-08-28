package vault

import (
	"context"
	"fmt"

	"google.golang.org/api/vault/v1"
)

const (
	ExportStatusUnknown ExportStatus = iota
	ExportStatusUnspecified
	ExportStatusCompleted
	ExportStatusFailed
	ExportStatusInProgress
)

//go:generate stringer -type=ExportStatus
type ExportStatus int

func (i ExportStatus) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

type CloudStorageFile struct {
	BucketName string
	ObjectName string
	Size       int64
	Md5Hash    string
}

type CloudStorageSink struct {
	Files []*CloudStorageFile
}

type ExportStatistics struct {
	ExportedArtifactCount int64
	TotalArtifactCount    int64
	SizeInBytes           int64
}

type Export struct {
	Id               string
	MatterId         string
	Name             string
	CreateTime       string
	Status           ExportStatus
	CloudStorageSink *CloudStorageSink
	Statistics       *ExportStatistics
}

func (v *VaultService) CreateDriveExportForMatter(ctx context.Context, matterId string, email string) (*Export, error) {
	exports, err := v.service.Matters.Exports.Create(matterId, &vault.Export{
		ExportOptions: &vault.ExportOptions{
			DriveOptions: &vault.DriveExportOptions{IncludeAccessInfo: false},
			Region:       "EUROPE",
		},
		Name: fmt.Sprintf("drive export of %q", email),
		Query: &vault.Query{
			AccountInfo:  &vault.AccountInfo{Emails: []string{email}},
			Corpus:       "DRIVE",
			DataScope:    "ALL_DATA",
			DriveOptions: &vault.DriveOptions{IncludeSharedDrives: false},
			SearchMethod: "ACCOUNT",
			Terms:        "",
		},
	}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return mapExport(exports), nil
}

func (v *VaultService) CreateEmailExportForMatter(ctx context.Context, matterId string, email string) (*Export, error) {
	exports, err := v.service.Matters.Exports.Create(matterId, &vault.Export{
		ExportOptions: &vault.ExportOptions{
			MailOptions: &vault.MailExportOptions{ExportFormat: "MBOX"},
			Region:      "EUROPE",
		},
		Name: fmt.Sprintf("email export of %q", email),
		Query: &vault.Query{
			AccountInfo:  &vault.AccountInfo{Emails: []string{email}},
			Corpus:       "MAIL",
			DataScope:    "ALL_DATA",
			MailOptions:  &vault.MailOptions{ExcludeDrafts: true},
			SearchMethod: "ACCOUNT",
			Terms:        "",
		},
	}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return mapExport(exports), nil
}

func (v *VaultService) FindExportsByMatterId(ctx context.Context, matterId string) ([]*Export, error) {
	exports, err := v.service.Matters.Exports.List(matterId).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return mapExports(exports.Exports), nil
}

func mapExports(e []*vault.Export) []*Export {
	exports := make([]*Export, len(e))
	for i, export := range e {
		exports[i] = mapExport(export)
	}
	return exports
}

func mapExport(e *vault.Export) *Export {
	return &Export{
		Id:               e.Id,
		MatterId:         e.MatterId,
		Name:             e.Name,
		CreateTime:       e.CreateTime,
		Status:           valueOf(e.Status),
		CloudStorageSink: mapCloudStorageSink(e.CloudStorageSink),
		Statistics:       mapExportStatistics(e.Stats),
	}
}

func mapExportStatistics(s *vault.ExportStats) *ExportStatistics {
	return &ExportStatistics{
		ExportedArtifactCount: s.ExportedArtifactCount,
		TotalArtifactCount:    s.TotalArtifactCount,
		SizeInBytes:           s.SizeInBytes,
	}
}

func mapCloudStorageSink(s *vault.CloudStorageSink) *CloudStorageSink {
	files := make([]*CloudStorageFile, 0, len(s.Files))
	for _, f := range s.Files {
		files = append(files, &CloudStorageFile{
			BucketName: f.BucketName,
			ObjectName: f.ObjectName,
			Size:       f.Size,
			Md5Hash:    f.Md5Hash,
		})
	}

	return &CloudStorageSink{Files: files}
}

func valueOf(s string) ExportStatus {
	switch s {

	case "EXPORT_STATUS_UNSPECIFIED":
		return ExportStatusUnspecified
	case "COMPLETED":
		return ExportStatusCompleted
	case "FAILED":
		return ExportStatusFailed
	case "IN_PROGRESS":
		return ExportStatusInProgress
	}
	return ExportStatusUnknown
}

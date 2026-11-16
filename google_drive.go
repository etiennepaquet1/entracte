package main

import (
	"bytes"
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const googleDriveFolderID = "1Z4OEP9Y9_0R4nTIX0qQD8PzzONrBSIbF"

func UploadReportToGoogleDrive(ctx context.Context, filename string, content []byte) (string, error) {
	service, err := drive.NewService(ctx, option.WithScopes(drive.DriveScope))
	if err != nil {
		return "", fmt.Errorf("create Google Drive service: %w", err)
	}

	file := &drive.File{
		Name:     filename,
		MimeType: "application/pdf",
		Parents:  []string{googleDriveFolderID},
	}
	uploaded, err := service.Files.Create(file).
		Media(bytes.NewReader(content)).
		Fields("id").
		Do()
	if err != nil {
		return "", fmt.Errorf("upload sales report to Google Drive: %w", err)
	}
	if uploaded.Id == "" {
		return "", fmt.Errorf("upload sales report to Google Drive: response did not include a file ID")
	}
	return uploaded.Id, nil
}

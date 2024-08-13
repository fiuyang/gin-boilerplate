package adapter

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type UploadBase64 struct {
	Folder string
	File   string
}

func (u *UploadBase64) GenerateFileName() (string, error) {
	// Extract the file extension from the base64 data URL
	dataParts := strings.Split(u.File, ",")
	if len(dataParts) != 2 {
		return "", errors.New("invalid base64 file data")
	}
	mimeType := strings.Split(dataParts[0], ";")[0]
	extension := getExtensionFromMimeType(mimeType)
	if extension == "" {
		return "", errors.New("file must have a valid extension")
	}
	guid := uuid.New()
	newFileName := fmt.Sprintf("%v/%v%v", u.Folder, guid.String(), extension)
	return newFileName, nil
}

func getExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "data:image/jpeg":
		return ".jpeg"
	case "data:image/jpg":
		return ".jpg"
	case "data:image/png":
		return ".png"
	default:
		return ""
	}
}

func (u *UploadBase64) GetFileContentType() (string, error) {
	dataParts := strings.Split(u.File, ",")
	if len(dataParts) != 2 {
		return "", errors.New("invalid base64 file data")
	}
	mimeType := strings.Split(dataParts[0], ";")[0]
	return mimeType, nil
}

func (u *UploadBase64) FileConvertToByteReader() (*bytes.Reader, error) {
	dataParts := strings.Split(u.File, ",")
	if len(dataParts) != 2 {
		return nil, errors.New("invalid base64 file data")
	}
	fileData := dataParts[1]

	// Decode the base64 data
	decodedData, err := base64.StdEncoding.DecodeString(fileData)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(decodedData), nil
}

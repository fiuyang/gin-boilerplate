package adapter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"path"
	"time"
)

type BucketList struct {
	Name      string
	CreatedAt time.Time
	Location  string
}

type UploadFile struct {
	Folder string
	File   *multipart.FileHeader
}

func (u *UploadFile) GenerateFileName() (string, error) {
	originFilename := u.File.Filename
	extension := path.Ext(originFilename)
	if extension == "" {
		return "", errors.New("file must be has extension")
	}
	guid := uuid.New()
	NewFileName := fmt.Sprintf("%v/%v%v", u.Folder, guid.String(), extension)
	return NewFileName, nil
}

func (u *UploadFile) GetFileContentType() string {
	fileheader := u.File.Header
	types, ok := fileheader["Content-Type"]
	var contentType string
	if ok {
		// This should be true!
		for _, x := range types {
			contentType = x
			// Most usually you will probably see only one
		}
	}
	return contentType
}

func (u *UploadFile) FileConvertToByteReader() (*bytes.Reader, error) {
	src, err := u.File.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	size := u.File.Size

	// Read the file into a byte slice
	bs := make([]byte, size)
	_, err = bufio.NewReader(src).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return nil, err
	}
	return bytes.NewReader(bs), nil
}

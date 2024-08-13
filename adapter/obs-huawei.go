package adapter

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"scylla/pkg/config"
)

type ObsAdapter interface {
	UploadFile(req *UploadFile) (fullUrl string, err error)
	UploadBase64(req *UploadBase64) (fullUrl string, err error)
	Delete(key string) error
}
type ObsAdapterImpl struct {
	Obsc        *obs.ObsClient
	FileBaseUrl string
	Bucket      string
}

func InitObsAdapter(conf config.ObsHuawei) (*ObsAdapterImpl, error) {
	endpoint := fmt.Sprintf("https://%v", conf.Endpoint)
	fileBaseUrl := fmt.Sprintf("https://%v.%v", conf.Bucket, conf.Endpoint)

	obsClient, err := obs.New(conf.Ak, conf.Sk, endpoint)
	if err != nil {
		fmt.Printf("Create obsClient error, errMsg: %s", err.Error())
		obsClient.Close()
		return nil, err
	}

	return &ObsAdapterImpl{Obsc: obsClient, FileBaseUrl: fileBaseUrl, Bucket: conf.Bucket}, nil
}

func (o *ObsAdapterImpl) UploadFile(req *UploadFile) (fullUrl string, err error) {
	key, err := req.GenerateFileName() // generate file name
	if err != nil {
		return
	}

	byteReader, err := req.FileConvertToByteReader() // convert form to byte reader
	if err != nil {
		return
	}

	input := &obs.PutObjectInput{}
	// Specify a bucket name.
	input.Bucket = o.Bucket
	// Specify the object (example/objectname as an example) to upload.
	input.Key = key
	input.ACL = obs.AclType(obs.AclPublicRead)
	input.ContentType = req.GetFileContentType()
	input.Body = byteReader
	// Upload you local file using streaming.
	_, err = o.Obsc.PutObject(input)
	if err == nil {
		fullUrl = fmt.Sprintf("%v/%v", o.FileBaseUrl, key)
		return
	}
	if obsError, ok := err.(obs.ObsError); ok {
		return fullUrl, obsError
	}
	return
}

func (o *ObsAdapterImpl) UploadBase64(req *UploadBase64) (fullUrl string, err error) {
	key, err := req.GenerateFileName() // generate file name
	if err != nil {
		return
	}

	byteReader, err := req.FileConvertToByteReader() // convert form to byte reader
	if err != nil {
		return
	}

	input := &obs.PutObjectInput{}
	// Specify a bucket name.
	input.Bucket = o.Bucket
	// Specify the object (example/objectname as an example) to upload.
	input.Key = key
	input.ACL = obs.AclType(obs.AclPublicRead)
	input.ContentType, _ = req.GetFileContentType()
	input.Body = byteReader
	// Upload you local file using streaming.
	_, err = o.Obsc.PutObject(input)
	if err == nil {
		fullUrl = fmt.Sprintf("%v/%v", o.FileBaseUrl, key)
		return
	}
	if obsError, ok := err.(obs.ObsError); ok {
		return fullUrl, obsError
	}
	return
}

func (o *ObsAdapterImpl) Delete(key string) error {
	input := &obs.DeleteObjectInput{
		Bucket: o.Bucket,
		Key:    key,
	}

	_, err := o.Obsc.DeleteObject(input)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			return fmt.Errorf("failed to delete file from OBS, errMsg: %s", obsError.Error())
		}
		return fmt.Errorf("failed to delete file from OBS, errMsg: %s", err.Error())
	}

	return nil
}

package services

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (videoUpload *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_KEY")+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(videoUpload.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (videoUpload *VideoUpload) loadPaths() error {
	err := filepath.Walk(videoUpload.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			videoUpload.Paths = append(videoUpload.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

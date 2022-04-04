package services

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}

func (videoUpload *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := videoUpload.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go videoUpload.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for x := 0; x < len(videoUpload.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}
	return nil
}

func (videoUpload *VideoUpload) uploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {
	for x := range in {
		err := videoUpload.UploadObject(videoUpload.Paths[x], uploadClient, ctx)
		if err != nil {
			videoUpload.Errors = append(videoUpload.Errors, videoUpload.Paths[x])
			log.Printf("error during the upload: %v, Error: %v", videoUpload.Paths[x], err)
			returnChan <- err.Error()
		}
		returnChan <- ""
	}
	returnChan <- "uploaded completed"
}

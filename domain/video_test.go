package domain_test

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"video-enconder-microsservice/domain"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIdIsNotAUuid(t *testing.T) {
	video := domain.NewVideo()

	video.ID = "ANY_ID"
	video.ResourceID = "ANY_ID"
	video.FilePath = "ANY_PATH"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.ResourceID = "ANY_ID"
	video.FilePath = "ANY_PATH"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Nil(t, err)
}

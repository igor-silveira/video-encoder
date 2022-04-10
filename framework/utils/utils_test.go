package utils_test

import (
	"github.com/stretchr/testify/require"
	"testing"
	"video-enconder-microsservice/framework/utils"
)

func TestIsJson(t *testing.T) {
	json := `{
  				"id": "525b5fd9-700d-4feb-89c0-415a1e6e148c",
  				"file_path": "convite.mp4",
  				"status": "pending"
			}`

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = `test`
	err = utils.IsJson(json)
	require.Error(t, err)
}

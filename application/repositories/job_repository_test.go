package repositories_test

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"video-enconder-microsservice/application/repositories"
	"video-enconder-microsservice/domain"
	"video-enconder-microsservice/framework/database"
)

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", "pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", "pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	job.Status = "complete"

	repoJob.Update(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.Status, job.Status)
}

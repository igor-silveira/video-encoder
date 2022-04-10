package services

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"os"
	"time"
	"video-enconder-microsservice/domain"
	"video-enconder-microsservice/framework/utils"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChan chan JobWorkerResult, jobService JobService, job domain.Job, workerID int) {

	for message := range messageChannel {
		err := utils.IsJson(string(message.Body))
		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = json.Unmarshal(message.Body, &jobService.VideoService.Video)
		jobService.VideoService.Video.ID = uuid.NewV4().String()
		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.Video.Validate()
		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.InsertVideo()
		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		job.Video = jobService.VideoService.Video
		job.OutputBucketPath = os.Getenv("OUTPUT_BUCKET_NAME")
		job.ID = uuid.NewV4().String()
		job.Status = "STARTING"
		job.CreatedAt = time.Now()

		_, err = jobService.JobRepository.Insert(&job)

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = &job
		err = jobService.Start()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChan <- returnJobResult(job, message, nil)
	}
}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	result := JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}
	return result
}

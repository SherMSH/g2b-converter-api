package jobs

import (
	"converterapi/internal/config"
	"converterapi/pkg/logger"
	"time"

	"github.com/go-co-op/gocron"
)

func Start() {
	logger.Infof("Launch the task scheduler...")
	params := config.Config.Jobs
	logger.Infof("Parameters: %+v", params)

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonMode()

	if params.ConvScanner.IsOn {
		if _, err := scheduler.Every(params.ConvScanner.Interval).Seconds().
			// StartAt(time.Now().Local().Add(time.Duration(params.ConvScanner.Interval) * time.Second)).
			Do(ConvScanner); err != nil {
			logger.Errorf("ConvScanner JOB err %v", err)
		}
	}
	scheduler.StartAsync()
}

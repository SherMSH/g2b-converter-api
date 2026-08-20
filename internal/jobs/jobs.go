package jobs

import (
	"converterapi/internal/config"
	d8procweb "converterapi/pkg/d8-proc-web"
	"converterapi/pkg/logger"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
)

func Start() {
	logger.Infof("Launch the task scheduler...")
	params := config.Config.Jobs
	logger.Infof("Parameters: %+v", params)

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonMode()

	if _, err := scheduler.Every(1).Hour().Do(signin); err != nil {
		logger.Errorf("Signin JOB err %v", err)
	}
	if params.ConvScanner.IsOn {
		if _, err := scheduler.Every(params.ConvScanner.Interval).Seconds().
			StartAt(time.Now().Local().Add(time.Duration(params.ConvScanner.Interval) * time.Second)).
			Do(ConvScanner); err != nil {
			logger.Errorf("ConvScanner JOB err %v", err)
		}
	}
	if params.TrnImporter.IsOn {
		// startAt := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 1, 0, 0, 0, time.Local)
		if _, err := scheduler.Every(params.TrnImporter.Interval).Hours().
			// StartAt(startAt).
			Do(TrnImporter); err != nil {
			logger.Errorf("TrnImporter JOB err %v", err)
		}
	}
	scheduler.StartAsync()
}

// signin поддерживает живую сессию с процессингом
func signin() {
	resp, err := d8procweb.Signin()
	if err != nil {
		logger.Errorf("[JOBS] Signin err: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("[JOBS] Signin status: %v", resp.Status)
	}
}

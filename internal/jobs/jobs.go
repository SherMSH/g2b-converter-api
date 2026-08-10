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

	if params.ConvScanner.IsOn {
		if _, err := scheduler.Every(params.ConvScanner.Interval).Seconds().
			// StartAt(time.Now().Local().Add(time.Duration(params.ConvScanner.Interval) * time.Second)).
			Do(ConvScanner); err != nil {
			logger.Errorf("ConvScanner JOB err %v", err)
		}
	}
	if _, err := scheduler.Every(30).Seconds().Do(signin); err != nil {
		logger.Errorf("Signin JOB err %v", err)
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

package d8procweb

import (
	"converterapi/internal/config"
	"converterapi/pkg/logger"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

var Client *http.Client

func Init() {
	jar, _ := cookiejar.New(nil)
	Client = &http.Client{
		Timeout: time.Duration(config.Config.App.ClientTimeoutSeconds) * time.Second,
		Jar:     jar,
	}
}

func Signin() (*http.Response, error) {
	if config.Config.Processing.Address == "" {
		return nil, fmt.Errorf("Portal URL not configured")
	}
	loginURL := config.Config.Processing.Address + "/api/login"

	authData := `{"login":"mustafokulovsh@humo.tj","password":"$B6d75737461666f6bbc458a4eebebe7d5226f297441bd01694cddf0b81008e585eec97f0b594a1078"}`

	if Client == nil {
		Init()
	}
	logger.Infof("trying %v... {timeout is %v seconds}", loginURL, Client.Timeout)
	resp, err := Client.Post(loginURL, "application/json", strings.NewReader(authData))
	if err != nil {
		logger.Errorf("trying %s error: %v", loginURL, err)
		return nil, fmt.Errorf("Portal connection failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp, nil
	}
	config.Config.Processing.Token = resp.Header["Set-Cookie"][0]
	return resp, nil
}

func Signout() (*http.Response, error) {
	if config.Config.Processing.Address == "" {
		return nil, fmt.Errorf("Portal URL not configured")
	}
	logoutURL := config.Config.Processing.Address + "/api/logout"

	logger.Infof("trying %v... {timeout is %v seconds}", logoutURL, Client.Timeout)
	if Client == nil {
		Init()
	}

	resp, err := Client.Post(logoutURL, "application/json", nil)
	if err != nil {
		logger.Errorf("trying %s error: %v", logoutURL, err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp, nil
	}

	config.Config.Processing.Token = ""
	return resp, nil
}

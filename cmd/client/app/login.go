package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func auth(client *resty.Client) (string, error) {
	if login == "" || password == "" {
		return "", errors.New("Please provide login and password flags or register first")
	}

	if serverURL == "" {
		return "", errors.New("Please specify server addr flag")
	}

	data := fmt.Sprintf("{\"login\": \"%s\", \"password\": \"%s\"}", login, password)

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(fmt.Sprintf("http://%s/api/user/login", serverURL))
	if err != nil {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Failed to login: %s\n", res.Body()))
	}

	return res.Header().Get("Authorization"), nil
}

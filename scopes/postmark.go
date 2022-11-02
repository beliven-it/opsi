package scopes

import (
	"encoding/json"
	"errors"
	"fmt"
	"opsi/helpers"
	"sync"
)

type Postmark struct {
	token        string
	slackWebhook string
	apiURL       string
}

func (p *Postmark) getGeneralPayload() postmarkEditRequest {
	return postmarkEditRequest{
		SmtpApiActivated:           true,
		RawEmailEnabled:            false,
		BounceHookUrl:              p.slackWebhook,
		PostFirstOpenOnly:          false,
		TrackOpens:                 true,
		TrackLinks:                 "HtmlAndText",
		IncludeBounceContentInHook: true,
		EnableSmtpApiErrorHooks:    false,
	}
}

func (p *Postmark) request(method string, endpoint string, body []byte, queryMap map[string]string) ([]byte, error) {
	return helpers.Request(method, p.apiURL+endpoint, body, queryMap, map[string]string{
		"Content-Type":             "application/json",
		"Accept":                   "application/json",
		"X-Postmark-Account-Token": p.token,
	})
}

func (p *Postmark) EditServer(serverID int) error {
	endpoint := fmt.Sprintf("/servers/%d", serverID)

	// Create the POST request payload
	payloadAsBytes, err := json.Marshal(p.getGeneralPayload())
	if err != nil {
		return err
	}

	// Execute the request
	_, err = p.request("PUT", endpoint, payloadAsBytes, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postmark) BulkEditServers() error {
	servers, err := p.GetServers()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	var errs = []error{}
	for _, server := range servers {
		wg.Add(1)

		go func(serverID int) {
			defer wg.Done()

			err := p.EditServer(serverID)
			if err != nil {
				errs = append(errs, err)
			}

		}(server.ID)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.New("something went wrong during the UPDATE")
	}

	return nil
}

func (p *Postmark) CreateServer(name string, color string) error {
	// Create the POST request payload
	payload := postmarkCreateRequest{
		Name:                name,
		Color:               color,
		postmarkEditRequest: p.getGeneralPayload(),
	}

	payloadAsBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Execute the request
	_, err = p.request("POST", "/servers", payloadAsBytes, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postmark) GetServers() ([]postmarkServer, error) {
	responseBody, err := p.request("GET", "/servers", nil, map[string]string{
		"count":  "100",
		"offset": "0",
	})

	if err != nil {
		return nil, err
	}

	var serverResponse postmarkServersResponse
	err = json.Unmarshal(responseBody, &serverResponse)
	if err != nil {
		return nil, err
	}
	return serverResponse.Servers, nil
}

func NewPostmark(apiURL string, token string, slackWebhook string) Postmark {
	return Postmark{
		apiURL:       apiURL,
		token:        token,
		slackWebhook: slackWebhook,
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/donpark/pam"
)

var transport *http.Transport
var baseURL string

type getPrompts struct {
	User  string `json:"user"`
	Token string `json:"token"`
}
type getPromptsResponse struct {
	Error   string
	Prompts []pam.Message
}

func getAuthPrompts(user, token string) (*getPromptsResponse, error) {
	url := baseURL + "/authPrompts"
	if isDebugMode {
		info("API-GETAUTH", fmt.Sprintf("Making request to %q", url))
	}

	b, err := json.Marshal(getPrompts{User: user, Token: token})
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)

	client := &http.Client{Transport: transport}
	resp, err := client.Post(url, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if isDebugMode {
		info("API-GETAUTH", fmt.Sprintf("Response: code=%d(%s),length=%d,content-type=%s", resp.StatusCode, resp.Status, resp.ContentLength, resp.Header.Get("Content-Type")))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	var response getPromptsResponse
	buf.Reset()
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return &response, nil
}

type authenticateRequest struct {
	User      string     `json:"user"`
	Token     string     `json:"token"`
	Responses [][]string `json:"responses"`
}

type authenticateResponse struct {
	Error   string
	Success bool
	Message string
}

func authenticate(user, token string, responses [][]string) (*authenticateResponse, error) {
	url := baseURL + "/authenticate"
	if isDebugMode {
		info("API-AUTHENTICATE", fmt.Sprintf("Making request to %q", url))
	}

	b, err := json.Marshal(authenticateRequest{User: user, Token: token, Responses: responses})
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)

	client := &http.Client{Transport: transport}
	resp, err := client.Post(url, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if isDebugMode {
		info("API-AUTHENTICATE", fmt.Sprintf("Response: code=%d(%s),length=%d,content-type=%s", resp.StatusCode, resp.Status, resp.ContentLength, resp.Header.Get("Content-Type")))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	var response authenticateResponse
	buf.Reset()
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return &response, nil
}

package util

import (
	"bytes"
	"fmt"
	"github.com/go-errors/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	request *http.Request
	client  *http.Client
}

func NewHttpClient() *HttpClient {
	httpClient := &HttpClient{
		client: &http.Client{},
	}
	return httpClient
}

func (h *HttpClient) MakeRequest(method, url string, data []byte) error {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	h.request = req
	return nil
}

func (h *HttpClient) SetHeaders(headers map[string]string) error {
	if h.request != nil {
		for k, v := range headers {
			h.request.Header.Add(k, v)
		}
		return nil
	}
	return fmt.Errorf("request is null! ")
}

func (h *HttpClient) SetBasicAuth(username string, password string) error {
	if h.request != nil {
		h.request.SetBasicAuth(username, password)
		return nil
	}
	return fmt.Errorf("request is null! ")
}

func (h *HttpClient) SetCookies(c *http.Cookie) error {
	if h.request != nil {
		h.request.AddCookie(c)
		return nil
	}
	return fmt.Errorf("request is null! ")
}

func (h *HttpClient) SetTimeout(timeout time.Duration) {
	h.client.Timeout = timeout
}

func (h *HttpClient) SetTransport(transport *http.Transport) {
	h.client.Transport = transport
}

func (h *HttpClient) DoRequest() (*http.Response, error) {
	if h.request == nil {
		return nil, errors.Errorf("make request first! ")
	}
	return h.client.Do(h.request)
}

func (h *HttpClient) GetResponseBody() ([]byte, error) {
	if h.request == nil {
		return nil, errors.Errorf("make request first! ")
	}
	resp, err := h.client.Do(h.request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		bRes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return bRes, nil
	}
	return nil, fmt.Errorf("response code:%v", resp.StatusCode)
}

package redmine

import "net/http"

type Client struct {
	endpoint string
	apikey   string
	*http.Client
}

func NewClient(endpoint, apikey string) *Client {
	return &Client{endpoint, apikey, http.DefaultClient}
}

type errorsResult struct {
	Errors []string `json:"errors"`
}

type IdName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Id struct {
	Id int `json:"id"`
}

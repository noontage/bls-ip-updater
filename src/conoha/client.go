package conoha

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	epIdentity   = "https://identity.tyo2.conoha.io"
	epDnsService = "https://dns-service.tyo2.conoha.io"
)

type Client struct {
	base *http.Client

	user     string
	password string
	tenantID string
}

func NewClient() *Client {
	return &Client{
		base:     &http.Client{},
		user:     os.Getenv("CNH_USER"),
		password: os.Getenv("CNH_PASSWORD"),
		tenantID: os.Getenv("CNH_TENANT_ID"),
	}
}

func (c Client) Test() error {
	token, err := c.getToken()
	if err != nil {
		return err
	}

	req, err := c.makeHttpRequest(http.MethodGet, epDnsService+"/v1/domains", token, []byte{})
	if err != nil {
		return err
	}

	resp, err := c.base.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	println(resp.StatusCode)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	println(string(b))

	return nil
}

func (c Client) makeHttpRequest(method, url, token string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c Client) getToken() (string, error) {
	idReq := map[string]interface{}{
		"auth": map[string]interface{}{
			"passwordCredentials": map[string]interface{}{
				"username": c.user,
				"password": c.password,
			},
			//"tenantId": c.tenantID,
		},
	}

	r, err := c.postJson(epIdentity+"/v2.0/tokens", idReq)
	if err != nil {
		return "", err
	}

	var resp struct {
		Access struct {
			Token struct {
				ID string `json:"id"`
			} `json:"token"`
		} `json:"access"`
	}

	if err := json.Unmarshal(r, &resp); err != nil {
		return "", err
	}

	return resp.Access.Token.ID, nil
}

func (c Client) postJson(url string, data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	resp, err := c.base.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return s, nil
}

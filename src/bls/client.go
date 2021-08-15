package bls

import (
	"bls-ip-updater/src/randstr"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

type Client struct {
	base   *http.Client
	secret string
	apiEndpoint string
}

func NewClient() *Client {
	return &Client{
		base:   &http.Client{},
		secret: os.Getenv("SECRET"),
		apiEndpoint: os.Getenv("BLS_ENDPOINT"),
	}
}

type cfReqest struct {
	Command   string `json:"command"`
	Salt      string `json:"salt"`
	Signature string `json:"signature"`
}

type cfResponse struct {
	Salt      string `json:"salt"`
	Signature string `json:"signature"`
	Data      string `json:"data"`
}

func (c Client) makeGetGlobalIPRequest() cfReqest {
	cmd := "global_ip"
	salt := randstr.NewString(8)
	return cfReqest{
		Command:   cmd,
		Salt:      salt,
		Signature: makeSignature(cmd, salt, c.secret),
	}
}

func (c Client) GetGlobalIP() (net.IP, error) {
	req := c.makeGetGlobalIPRequest()

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.base.Post(c.apiEndpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status not ok statusCode=%d", resp.StatusCode)
	}

	var result cfResponse
	if err := json.Unmarshal(j, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	if makeSignature(result.Data, result.Salt, c.secret) != result.Signature {
		return nil, errors.New("invalid signature")
	}

	return net.ParseIP(result.Data), nil
}

func makeSignature(data, salt, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data + "." + salt))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SmarterMailClient struct {
	httpClient        *http.Client
	SmarterMailConfig SmarterMailConfigDTO
	ApiBaseUrl        string
	ApiTokenAuth      string
}

func InitSmarterMail(SmarterMailConfig SmarterMailConfigDTO) (*SmarterMailClient, error) {

	c := &SmarterMailClient{
		httpClient:        &http.Client{},
		SmarterMailConfig: SmarterMailConfig,
		ApiBaseUrl:        fmt.Sprintf("https://%v/api/v1", SmarterMailConfig.Host),
	}

	tokenApi, err := c.Authenticate(c.SmarterMailConfig.Username, c.SmarterMailConfig.Password)
	if err != nil {
		return nil, err
		// log.Fatalln("Erro ao realizar a autenticação com o usuário", c.SmarterMailConfig.Username, "!!")
	}
	c.ApiTokenAuth = tokenApi

	return c, nil
}

func (c *SmarterMailClient) makeRequest(method, ApiUrlRoute string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.ApiBaseUrl+ApiUrlRoute, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.ApiTokenAuth))

	if err != nil {
		return nil, err
	}

	return req, nil

}

func (c *SmarterMailClient) Get(ApiUrlRoute string) (*http.Response, error) {

	req, err := c.makeRequest("GET", ApiUrlRoute, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil

}
func (c *SmarterMailClient) Post(ApiUrlRoute string, payload io.Reader, customHeader map[string]string) (*http.Response, error) {

	req, err := c.makeRequest("POST", ApiUrlRoute, payload)
	if err != nil {
		fmt.Println(err)
	}

	for i, v := range customHeader {
		req.Header.Set(i, v)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (c *SmarterMailClient) Authenticate(username string, password string) (string, error) {

	CredencialsInputDTO := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	}

	CredencialsJsonPayload, err := json.Marshal(CredencialsInputDTO)
	if err != nil {
		return "", err
	}
	CredencialsPayloadBuf := bytes.NewBuffer(CredencialsJsonPayload)

	resp, err := c.Post("/auth/authenticate-user", CredencialsPayloadBuf, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ResponseBodyBytes, _ := io.ReadAll(resp.Body)

	ResponseBody := struct {
		Success     bool   `json:"success"`
		AccessToken string `json:"accessToken"`
	}{}

	err = json.Unmarshal(ResponseBodyBytes, &ResponseBody)
	if err != nil {
		return "", err
	}

	if !ResponseBody.Success {
		return "", fmt.Errorf("erro ao autenticar com o usuário %v", username)
	}

	return ResponseBody.AccessToken, nil

}

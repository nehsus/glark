package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//Lark provides methods to interact with lark using lark APIs
type Lark interface {
	GetAccessToken(appDetails AppDetails) (LarkToken, error)
	GetBotGroups(token string) (LarkData, error)
	SendMessage(token string, larkMessageRequest LarkMessageRequest) int
}

type larkRequest struct {
	client *http.Client
}

var lark Lark

func (lr *larkRequest) GetAccessToken(appDetails AppDetails) (LarkToken, error) {
	var token LarkToken
	payload, err := json.Marshal(appDetails)
	if err != nil {
		return token, err
	}
	req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/auth/v3/app_access_token/internal", bytes.NewBuffer(payload))
	if err != nil {
		return token, err
	}
	resp, err := lr.client.Do(req)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (lr *larkRequest) GetBotGroups(token string) (LarkData, error) {
	var larkData LarkData
	req, err := http.NewRequest("GET", "https://open.larksuite.com/open-apis/chat/v4/list", nil)
	if err != nil {
		return larkData, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := lr.client.Do(req)
	if err != nil {
		return larkData, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return larkData, err
	}
	err = json.Unmarshal([]byte(body), &larkData)
	if err != nil {
		return larkData, err
	}
	return larkData, nil
}

func (lr *larkRequest) SendMessage(token string, larkMessageRequest LarkMessageRequest) int {
	payload, err := json.Marshal(larkMessageRequest)
	if err != nil {
		return http.StatusBadRequest
	}
	req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/message/v4/send/", bytes.NewBuffer(payload))
	req.Header.Add("Authorization", "Bearer "+token)
	resp, _ := lr.client.Do(req)
	return resp.StatusCode
}

//InitLark is used to initialise lark
func InitLark(l Lark) {
	lark = l
}

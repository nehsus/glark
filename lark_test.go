package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LarkSuite struct {
	suite.Suite
	lark   *larkRequest
	client *http.Client
}

func (l *LarkSuite) SetupSuite() {
	client := http.Client{}
	// InitLark(&larkRequest{client: &client})
	l.lark = &larkRequest{client: &client}
	l.client = &client
	appID = os.Getenv("TEST_APP_ID")
	appSecret = os.Getenv("TEST_APP_SECRET")
	testChatID = os.Getenv("TEST_CHAT_ID")
}

func (l *LarkSuite) TearDownSuite() {
	l.client.CloseIdleConnections()
}

func TestStoreSuite(t *testing.T) {
	l := new(LarkSuite)
	suite.Run(t, l)
}

func (l *LarkSuite) TestGetAccessToken() {
	appDetails := AppDetails{AppID: appID, AppSecret: appSecret}
	larkToken, err := l.lark.GetAccessToken(appDetails)
	if err != nil {
		l.T().Fatal(err)
	}
	if larkToken.Msg != "ok" {
		l.T().Errorf("Incorrect message, wanted %v, got %v", "ok", larkToken.Msg)
	}
}

func (l *LarkSuite) TestGetBotGroups() {
	var token LarkToken
	appDetails := AppDetails{AppID: appID, AppSecret: appSecret}

	payload, err := json.Marshal(appDetails)
	if err != nil {
		l.T().Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/auth/v3/app_access_token/internal", bytes.NewBuffer(payload))
	if err != nil {
		l.T().Fatal(err)
	}
	resp, err := l.client.Do(req)
	if err != nil {
		l.T().Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.T().Fatal(err)
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		l.T().Fatal(err)
	}
	_, err = l.lark.GetBotGroups(token.TenantAccessToken)
	if err != nil {
		l.T().Fatal(err)
	}
}

func (l *LarkSuite) TestSendMessage() {
	var token LarkToken
	appDetails := AppDetails{AppID: appID, AppSecret: appSecret}

	payload, err := json.Marshal(appDetails)
	if err != nil {
		l.T().Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/auth/v3/app_access_token/internal", bytes.NewBuffer(payload))
	if err != nil {
		l.T().Fatal(err)
	}
	resp, err := l.client.Do(req)
	if err != nil {
		l.T().Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.T().Fatal(err)
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		l.T().Fatal(err)
	}
	var larkMessageRequest = LarkMessageRequest{
		ChatID:  testChatID,
		MsgType: "text",
		Content: Content{
			Text: "This is test text from test bot!!",
		},
	}
	statusCode := l.lark.SendMessage(token.TenantAccessToken, larkMessageRequest)
	if statusCode != http.StatusOK {
		l.T().Errorf("Incorrect status code, wanted %v, got %v", http.StatusOK, statusCode)
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleRequest(t *testing.T) {
	mockLark := InitMockLark()

	var request events.APIGatewayProxyRequest
	var payload Payload
	payload.DashboardID = 123
	payload.Message = "Test Message"
	payload.EvalMatches = append(payload.EvalMatches, EvalMatches{Metric: " Metric"})
	bytes, _ := json.Marshal(payload)
	request.Body = string(bytes)

	appDetails := AppDetails{AppID: appID, AppSecret: appSecret}
	larkMessageRequest := LarkMessageRequest{MsgType: "text", ChatID: "test_chat_id", Content: Content{Text: "Test Message Metric"}}
	mockLark.On("GetAccessToken", appDetails).Return(LarkToken{TenantAccessToken: "token"}, nil)
	var groups []Group
	groups = append(groups, Group{ChatID: "test_chat_id", Name: grafanaChatName})
	mockLark.On("GetBotGroups", "token").Return(LarkData{Data: Data{Groups: groups}}, nil)
	mockLark.On("SendMessage", "token", larkMessageRequest).Return(http.StatusOK)

	response, err := handleRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Incorrect status code, wanted %v, got %v", http.StatusOK, response.StatusCode)
	}
	mockLark.AssertExpectations(t)
}

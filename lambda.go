package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("%v\n", request)
	appDetails := AppDetails{
		AppID:     appID,
		AppSecret: appSecret,
	}
	token, err := lark.GetAccessToken(appDetails)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}
	larkData, err := lark.GetBotGroups(token.TenantAccessToken)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, err
	}
	var chatID string
	for _, group := range larkData.Data.Groups {
		if group.Name == grafanaChatName {
			//Group found
			chatID = group.ChatID
		}
	}

	var payload Payload
	err = json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, err
	}
	fmt.Println("Here", payload)
	var larkMessageRequest = LarkMessageRequest{
		ChatID:  chatID,
		MsgType: "text",
		Content: Content{
			Text: payload.Message + payload.EvalMatches[0].Metric,
		},
	}
	statusCode := lark.SendMessage(token.TenantAccessToken, larkMessageRequest)
	if statusCode != 200 {
		return events.APIGatewayProxyResponse{Body: "Bad Request", StatusCode: statusCode}, errors.New("Bad request")
	}
	return events.APIGatewayProxyResponse{Body: "Message Sent", StatusCode: statusCode}, nil

}

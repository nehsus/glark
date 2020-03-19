package glark

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// AppDetails is a struct to store credentials from our Lark application
type AppDetails struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

/*
LarkToken is a struct to store credentials obtained from Lark
Tokens required to access server-side APIs:
  - app_access_token: To access APIs with app resources.
  - tenant_access_token: To access APIs with company resources.
*/
type LarkToken struct {
	AppAccessToken    string `json:"app_access_token"`
	Code              int    `json:"code"`
	Expire            int    `json:"expire"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}

type LarkData struct {
	Code int    `json:"code"`
	Data Data   `json:"data"`
	Msg  string `json:"msg"`
}

type Data struct {
	Groups    []Groups `json:"groups"`
	HasMore   bool     `json:"has_more"`
	PageToken string   `json:"page_token"`
}

type Groups struct {
	Avatar      string `json:"avatar"`
	ChatID      string `json:"chat_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	OwnerOpenID string `json:"owner_open_id"`
	OwnerUserID string `json:"owner_user_id"`
}

// Payload is a struct to store details from Grafana
type Payload struct {
	DashboardID int           `json:"dashboardId"`
	EvalMatches []EvalMatches `json:"evalMatches"`
	ImageURL    string        `json:"imageUrl"`
	Message     string        `json:"message"`
	OrgID       int           `json:"orgId"`
	PanelID     int           `json:"panelId"`
	RuleID      int           `json:"ruleId"`
	RuleName    string        `json:"ruleName"`
	RuleURL     string        `json:"ruleUrl"`
	State       string        `json:"state"`
	Tags        Tags          `json:"tags"`
	Title       string        `json:"title"`
}

type EvalMatches struct {
	Value  int    `json:"value"`
	Metric string `json:"metric"`
	Tags   Tags   `json:"tags"`
}

type Tags struct {
	TagName string `json:"tag name"`
}

//LarkMessageRequest is a struct to store the body for Lark message
type LarkMessageRequest struct {
	MsgType string  `json:"msg_type"`
	ChatID  string  `json:"chat_id"`
	Content Content `json:"content"`
}
type Content struct {
	Text string `json:"text"`
}

// Handler function Using AWS Lambda Proxy Request
func glark(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//Get environment variabls from AWS Lambda
	AppID := os.Getenv("app_id")
	AppSecret := os.Getenv("app_secret")
	ChatName := os.Getenv("chat_name")

	// Payload will be used to take the json Payload from Grafana
	var (
		payload   Payload
		larkToken LarkToken
		larkData  LarkData
		chatID    string
	)

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Unmarshal the json, return 404 if error.
	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	larkPayload := AppDetails{
		AppID:     AppID,
		AppSecret: AppSecret,
	}

	// Marshal the application credentials so we can request a token from Lark
	larkBytesRepresentation, err := json.Marshal(larkPayload)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	// Lark request
	req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/auth/v3/app_access_token/internal", bytes.NewBuffer(larkBytesRepresentation))

	// Send HTTP request to Lark
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(body)))

	// Unmarshal the JSON to store the LarkToken
	err = json.Unmarshal([]byte(body), &larkToken)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	bearer := "Bearer " + larkToken.TenantAccessToken // larkToken.AppAccessToken

	//New Request to obtain groups the bot is part of
	req, err = http.NewRequest("GET", "https://open.larksuite.com/open-apis/chat/v4/list", nil)

	// Add authorization headers
	req.Header.Add("Authorization", bearer)

	// Send HTTP request to Lark
	resp, err = client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(body)))

	// Unmarshal the JSON to store the LarkToken
	err = json.Unmarshal([]byte(body), &larkData)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
	}

	var groupSlice []Groups
	//Need to search "chat_id" in LarkData for the ChatID supplied by env
	for i := range groupSlice {
		if groupSlice[i].Name == ChatName {
			//Group found
			chatID = larkData.Data.Groups[i].ChatID
		}
	}

	// Then we prepare the message payload
	reqPayload := LarkMessageRequest{
		ChatID:  chatID,
		MsgType: "text",
		Content: Content{
			Text: payload.Message + payload.EvalMatches[0].Metric,
		},
	}

	bytesRepresentation, err := json.Marshal(reqPayload)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err

	}

	// Create a new Lark request
	req, err = http.NewRequest("POST", "https://open.larksuite.com/open-apis/message/v4/send/", bytes.NewBuffer(bytesRepresentation))

	// Add authorization headers
	req.Header.Add("Authorization", bearer)

	// Send HTTP request to Lark
	resp, err = client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(body)))

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	// Returning response with Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: "Done", StatusCode: 200}, nil
}

func main() {
	lambda.Start(glark)
}

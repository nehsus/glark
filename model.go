package main

var appID, appSecret, testChatID, grafanaChatName string

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

//LarkData is a struct to store group information for Lark's Bot
type LarkData struct {
	Code int    `json:"code"`
	Data Data   `json:"data"`
	Msg  string `json:"msg"`
}

//Data is supporting struct used to store group details
type Data struct {
	Groups    []Group `json:"groups"`
	HasMore   bool    `json:"has_more"`
	PageToken string  `json:"page_token"`
}

//Group is a supporting struct used to store group details
type Group struct {
	Avatar      string `json:"avatar"`
	ChatID      string `json:"chat_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	OwnerOpenID string `json:"owner_open_id"`
	OwnerUserID string `json:"owner_user_id"`
}

//Payload is a struct to store details from Grafana
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

//EvalMatches is a supporting struct used to store metrics, tags and values from Grafana
type EvalMatches struct {
	Value  int    `json:"value"`
	Metric string `json:"metric"`
	Tags   Tags   `json:"tags"`
}

//Tags is supporting struct used to store grafana tags
type Tags struct {
	TagName string `json:"tag name"`
}

//LarkMessageRequest is a struct to store the body for a Lark message
type LarkMessageRequest struct {
	MsgType string  `json:"msg_type"`
	ChatID  string  `json:"chat_id"`
	Content Content `json:"content"`
}

//Content is a supporting struct to store message content
type Content struct {
	Text string `json:"text"`
}

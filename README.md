[![Build Status](https://travis-ci.org/nehsus/glark.svg?branch=master)](https://travis-ci.org/nehsus/glark)
<br />
# Glark
Glark is an alertManager for Grafana with Lark, written in Go

## Installation
```bash
go get github.com/nehsus/glark
```
## Working Steps

 * Setup [AWS Lambda](https://aws.amazon.com/lambda/), configured with Go

 * Setup a custom notification channel in Grafana with a webhook to Lambda

 * Create a new [Lark Bot](https://open.larksuite.com/document/uMDNxEjLzQTMx4yM0ETM/uUjMyEjL1IjMx4SNyITM) and obtain credentials: 
    - AppID
    - AppSecret

 * Obtain the [API access token](https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/ukjMyEjL5IjMx4SOyITM) from Lark

 * Obtain a new bearertoken from Lark with the application credentials \
Request method: POST \
Request address: https://open.larksuite.com/open-apis/auth/v3/app_access_token/internal \
Request header: \
key	value \
Content-Type	application/json \
Request example: 
```json
{ 
    "app_id": "supersecretid", 
    "app_secret": "supersecretsecret" 
}
```

 * Obtain All [Groups](https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uYjMxUjL2ITM14iNyETN) to which a user belongs and note down the chat_id

 * Configure Lambda with environment variables:
    - app_id
    - app_secret
    - chat_id

 * Invoke the message sending API: \
Request method: POST \
Request address: https://open.larksuite.com/open-apis/message/v4/send/ \
Request header: \
key	value \
Authorization	Bearer tenant_access_token \
Content-Type	application/json \
Request example:
```json
{
   "chat_id":"oc_xxx", 
   "email":"test@gmail.com", 
    "msg_type":"text",
    "content":{
        "text":"test notification"
    }
}
```

## Pending Features
This project is currently under construction. Issues presently being worked on:
 
 * Only one bot and chat-group is implemented
 * You will need to obtain the chat_id of the group by this [method](https://open.larksuite.com/document/uMzMyEjLzMjMx4yMzITM/uETM1EjLxETNx4SMxUTM) (Currently under const.)

## Contributing
This package was developed in my free time. Contributions from everone are welcome to make this a more wholesome and streamlined experience. If you find any bugs or think there should be a particular feature included, feel free to open up a new issue or pull request.

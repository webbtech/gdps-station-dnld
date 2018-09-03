package main

import (
	"context"
	"encoding/json"

	"github.com/pulpfree/gdps-fs-dwnld/config"
	"github.com/pulpfree/gdps-fs-dwnld/fuelsale"
	"github.com/pulpfree/gdps-fs-dwnld/model"
	"github.com/pulpfree/gdps-fs-dwnld/validate"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var cfg *config.Config

const defaultsFilePath = "./defaults.yaml"

func init() {
	cfg = &config.Config{
		DefaultsFilePath: defaultsFilePath,
	}
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	var err error
	var eRes string

	// If this is a ping test, intercept and return
	if req.HTTPMethod == "GET" {
		log.Info("Ping test in handleRequest")
		return events.APIGatewayProxyResponse{Body: "pong", Headers: hdrs, StatusCode: 200}, nil
	}

	// Check for auth header
	/*if req.Headers["Authorization"] == "" {
		eRes = setErrorResponse(401, "Unauthorized", "Missing Authorization header")
		return events.APIGatewayProxyResponse{Body: eRes, Headers: hdrs, StatusCode: 401}, nil
	}

	// Set auth config
	auth, err := auth.New(&auth.Config{
		ClientID:       cfg.CognitoClientID,
		PoolID:         cfg.CognitoPoolID,
		Region:         cfg.CognitoRegion,
		JwtAccessToken: req.Headers["Authorization"],
	})
	if err != nil {
		eRes = setErrorResponse(500, "Authentication", err.Error())
		return events.APIGatewayProxyResponse{Body: eRes, Headers: hdrs, StatusCode: 500}, nil
	}

	// Validate JWT Token
	err = auth.Validate()
	if err != nil {
		eRes = setErrorResponse(401, "Authentication", err.Error())
		return events.APIGatewayProxyResponse{Body: eRes, Headers: hdrs, StatusCode: 401}, nil
	}*/

	// Set and validate request params
	var r *model.RequestInput
	json.Unmarshal([]byte(req.Body), &r)
	reqVars, err := validate.RequestInput(r)
	if err != nil {
		eRes = setErrorResponse(500, "RequestValidation", err.Error())
		return events.APIGatewayProxyResponse{Body: eRes, Headers: hdrs, StatusCode: 500}, nil
	}

	// Process request
	report, err := fuelsale.New(reqVars, cfg)
	if err != nil {
		eRes = setErrorResponse(500, "RequestValidation", err.Error())
		return events.APIGatewayProxyResponse{Body: eRes, Headers: hdrs, StatusCode: 500}, nil
	}

	var url string
	err = report.Create()
	if err != nil {
		eRes = setErrorResponse(500, "ProcessError", err.Error())
		return events.APIGatewayProxyResponse{Body: eRes, Headers: hdrs, StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: url, Headers: hdrs, StatusCode: 201}, nil
}

func main() {
	lambda.Start(handleRequest)
}

// ======================== Helper Function ================================= //

func setErrorResponse(status int, errType, message string) string {

	err := model.ErrorResponse{
		Status:  status,
		Type:    errType,
		Message: message,
	}
	log.Errorf("Error: status: %d, type: %s, message: %s", err.Status, err.Type, err.Message)
	res, _ := json.Marshal(&err)

	return string(res)
}

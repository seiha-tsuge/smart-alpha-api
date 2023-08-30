package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

func HandleRequest(ctx context.Context) (Response, error) {
	env := os.Getenv("APP_ENV")

	message := "Hello, this is the default environment."
	if env != "" {
		message = fmt.Sprintf("Hello, this is the %s environment.", env)
	}

	return Response{
		StatusCode: 200,
		Body:       message,
	}, nil
}

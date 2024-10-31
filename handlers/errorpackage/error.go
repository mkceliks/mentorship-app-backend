package errorPackage

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrNoSuchKey   = errors.New("NoSuchKey")
)

func ServerError(message string) (events.APIGatewayProxyResponse, error) {
	log.Printf(message)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       message,
		Headers:    wrapper.SetHeadersGet("text/plain"),
	}, errors.New(message)
}

func ClientError(status int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       message,
		Headers:    wrapper.SetHeadersGet("text/plain"),
	}, errors.New(message)
}

func HandleS3Error(err error) (events.APIGatewayProxyResponse, error) {
	switch {
	case errors.Is(err, ErrNoSuchKey):
		return ClientError(http.StatusNotFound, "File not found")
	default:
		log.Printf("S3 error: %v", err)
		return ServerError("Internal server error")
	}
}

package errorpackage

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"mentorship-app-backend/handlers/wrapper"
	"net/http"
)

var (
	ErrNoSuchKey = errors.New("NoSuchKey")
)

func ServerError(message string) (events.APIGatewayProxyResponse, error) {
	err := fmt.Errorf("server error: %s", message)
	log.Println(err)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       fmt.Sprintf(`{"error": "%s"}`, message),
		Headers:    wrapper.SetHeadersGet("application/json"),
	}, err
}

func ClientError(status int, message string) (events.APIGatewayProxyResponse, error) {
	err := fmt.Errorf("client error: %s", message)
	log.Println(err)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       fmt.Sprintf(`{"error": "%s"}`, message),
		Headers:    wrapper.SetHeadersGet("application/json"),
	}, err
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

func NewErrorResponse(status int, err error) (events.APIGatewayProxyResponse, error) {
	log.Printf("Error: %v", err)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		Headers:    wrapper.SetHeadersGet("application/json"),
	}, err
}

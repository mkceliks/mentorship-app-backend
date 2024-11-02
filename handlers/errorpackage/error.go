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
	ErrNoSuchKey          = errors.New("NoSuchKey")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func HandleS3Error(err error) (events.APIGatewayProxyResponse, error) {
	switch {
	case errors.Is(err, ErrNoSuchKey):
		return ClientError(http.StatusNotFound, "File not found")
	default:
		log.Printf("S3 error: %v", err)
		return ServerError("Internal server error")
	}
}

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

func IsUserAlreadyExistsError(err error) bool {
	return errors.Is(err, ErrUserAlreadyExists)
}

func IsInvalidCredentialsError(err error) bool {
	return errors.Is(err, ErrInvalidCredentials)
}

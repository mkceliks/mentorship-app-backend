package errorpackage

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"strings"
)

var (
	ErrNoSuchKey            = errors.New("NoSuchKey")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrMissingAuthorization = errors.New("authorization header is missing")
	ErrMissingToken         = errors.New("ID token is missing")
	ErrInvalidTokenFormat   = errors.New("invalid ID token format")
	ErrEmailNotFound        = errors.New("email not found in ID token")
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
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, err
}

func ClientError(status int, message string) (events.APIGatewayProxyResponse, error) {
	err := fmt.Errorf("client error: %s", message)
	log.Println(err)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       fmt.Sprintf(`{"error": "%s"}`, message),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, err
}

func IsInvalidConfirmationCodeError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "InvalidParameterException") &&
		strings.Contains(err.Error(), "Invalid code")
}

func IsExpiredConfirmationCodeError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "ExpiredCodeException")
}

func IsUserAlreadyExistsError(err error) bool {
	return errors.Is(err, ErrUserAlreadyExists)
}

func IsInvalidCredentialsError(err error) bool {
	return errors.Is(err, ErrInvalidCredentials)
}

func IsDynamoDBNotFoundError(err error) bool {
	return errors.Is(err, ErrNoSuchKey)
}

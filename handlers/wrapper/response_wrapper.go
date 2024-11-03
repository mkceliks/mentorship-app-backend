package wrapper

import "mentorship-app-backend/entity"

func SetAccessControl() map[string]string {
	return entity.AccessControl{
		"Access-Control-Allow-Origin": "*",
	}
}

func SetHeadersDelete() map[string]string {
	return entity.Headers{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}
}

func SetHeadersGet(contentType string) map[string]string {
	if contentType == "" {
		contentType = "application/json"
	}
	return entity.Headers{
		"Content-Type":                 contentType,
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}
}

func SetHeadersPost() map[string]string {
	return entity.Headers{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, x-file-content-type",
	}
}

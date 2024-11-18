package entity

type IDTokenPayload struct {
	Email         string `json:"email"`
	CustomRole    string `json:"custom:role"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
	Sub           string `json:"sub"`
}

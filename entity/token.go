package entity

type IDTokenPayload struct {
	Email      string `json:"email"`
	CustomRole string `json:"custom:role"`
	Name       string `json:"name"`
	Sub        string `json:"sub"`
}

package entity

type ConfirmRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

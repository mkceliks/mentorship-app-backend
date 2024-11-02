package entity

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

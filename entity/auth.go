package entity

type AuthRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	ProfilePictureBase64 string `json:"profile_picture_base64"`
}

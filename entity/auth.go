package entity

type AuthRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Role           string `json:"role"`
	ProfilePicture string `json:"profile_picture"`
	FileName       string `json:"file_name"`
}

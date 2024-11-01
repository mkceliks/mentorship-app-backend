package entity

type File struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}

type UploadRequest struct {
	Filename    string `json:"file_name"`
	FileContent string `json:"file_content"`
}

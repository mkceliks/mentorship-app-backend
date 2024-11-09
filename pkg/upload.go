package pkg

import (
	"fmt"
	"mentorship-app-backend/entity"
)

func (c *Client) UploadProfilePicture(fileName, base64Image, contentType string) (*entity.UploadResponse, error) {
	uploadReq := entity.UploadRequest{
		Filename:    fileName,
		FileContent: base64Image,
	}

	var uploadResp entity.UploadResponse
	resp, err := c.client.R().
		SetHeader("x-file-content-type", contentType).
		SetBody(uploadReq).
		SetResult(&uploadResp).
		Post("/upload")

	if err != nil {
		return nil, fmt.Errorf("failed to call upload API: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("upload API returned status %d: %s", resp.StatusCode(), resp.String())
	}

	return &uploadResp, nil
}

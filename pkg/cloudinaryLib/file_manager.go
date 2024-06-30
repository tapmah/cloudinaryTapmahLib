package cloudinarylib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type CloudinaryResource struct {
	AssetID       string  `json:"asset_id"`
	PublicID      string  `json:"public_id"`
	AssetFolder   string  `json:"asset_folder"`
	Filename      string  `json:"filename"`
	DisplayName   string  `json:"display_name"`
	Format        string  `json:"format"`
	Version       int     `json:"version"`
	ResourceType  string  `json:"resource_type"`
	Type          string  `json:"type"`
	CreatedAt     string  `json:"created_at"`
	UploadedAt    string  `json:"uploaded_at"`
	Bytes         int     `json:"bytes"`
	BackupBytes   int     `json:"backup_bytes"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	AspectRatio   float64 `json:"aspect_ratio"`
	Pixels        int     `json:"pixels"`
	URL           string  `json:"url"`
	SecureURL     string  `json:"secure_url"`
	Status        string  `json:"status"`
	AccessMode    string  `json:"access_mode"`
	AccessControl string  `json:"access_control"`
	Etag          string  `json:"etag"`
	CreatedBy     struct {
		AccessKey string `json:"access_key"`
	} `json:"created_by"`
	UploadedBy struct {
		AccessKey string `json:"access_key"`
	} `json:"uploaded_by"`
}

type CloudinaryResponse struct {
	TotalCount int                  `json:"total_count"`
	Time       int                  `json:"time"`
	Resources  []CloudinaryResource `json:"resources"`
}

func (b *BaseCloudinaryStuct) GetAllResources(maxFiles int) ([]CloudinaryResource, error) {
	params := url.Values{}
	params.Set("type", "upload")
	params.Set("prefix", "")
	params.Set("max_results", fmt.Sprint(maxFiles))
	params.Set("max_files", fmt.Sprint(maxFiles))

	apiURL := "https://api.cloudinary.com/v1_1/" + b.CloudName + "/resources/search"
	url := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []CloudinaryResource{}, err
	}

	req.Header.Set("Authorization", "Basic "+b.Base64Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []CloudinaryResource{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []CloudinaryResource{}, err
	}

	var response CloudinaryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []CloudinaryResource{}, err
	}
	ans := make([]CloudinaryResource, 0, len(response.Resources))
	ans = append(ans, response.Resources...)
	return ans, nil
}

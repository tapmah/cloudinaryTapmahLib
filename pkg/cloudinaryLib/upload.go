package cloudinarylib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"time"
)

type UploadResponseData struct {
	AssetID           string        `json:"asset_id"`
	PublicID          string        `json:"public_id"`
	Version           int           `json:"version"`
	VersionID         string        `json:"version_id"`
	Signature         string        `json:"signature"`
	Width             int           `json:"width"`
	Height            int           `json:"height"`
	Format            string        `json:"format"`
	ResourceType      string        `json:"resource_type"`
	CreatedAt         string        `json:"created_at"`
	Tags              []interface{} `json:"tags"`
	Pages             int           `json:"pages"`
	Bytes             int           `json:"bytes"`
	Type              string        `json:"type"`
	Etag              string        `json:"etag"`
	Placeholder       bool          `json:"placeholder"`
	URL               string        `json:"url"`
	SecureURL         string        `json:"secure_url"`
	AssetFolder       string        `json:"asset_folder"`
	DisplayName       string        `json:"display_name"`
	Overwritten       bool          `json:"overwritten"`
	OriginalFilename  string        `json:"original_filename"`
	OriginalExtension string        `json:"original_extension"`
	APIKey            string        `json:"api_key"`
	Error             struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (b BaseCloudinaryStuct) UploadFile(filePath string) (UploadResponseData, error) {

	url := "https://api.cloudinary.com/v1_1/" + b.CloudName + "/image/upload"

	var res UploadResponseData // Create Upload Response

	bodyBuf := &bytes.Buffer{}                 // Create new multipart buffer
	bodyWriter := multipart.NewWriter(bodyBuf) // Create buff writer

	file, err := os.Open(filePath) // Open file by file Path, return link to the file
	if err != nil {
		return res, err
	}
	defer file.Close()
	fileWriter, err := bodyWriter.CreateFormFile("file", filePath) // Open file by file Path for writer
	if err != nil {
		return res, err
	}
	_, err = io.Copy(fileWriter, file) // Copy bytes from file to buffer
	if err != nil {
		return res, err
	}

	params := map[string]string{} // Required params for signature

	b.TimeStamp = fmt.Sprintf("%d", time.Now().Unix()) // Create variable with time.Now

	value := reflect.ValueOf(b) // Convert UploadAdditionalParams to reflect.Value

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i).String()
		tag := value.Type().Field(i).Tag.Get("json")
		if len(field) > 0 && len(tag) > 0 {
			params[tag] = field                     // Add params
			err = bodyWriter.WriteField(tag, field) // Add request params
			if err != nil {
				return res, err
			}
		}
	}

	sortedParams := sortParams(params) + b.ApiSecret // Create string with required data
	signature := ComputeSHA1(sortedParams)           // Calculate SHA1 Sum

	err = bodyWriter.WriteField("signature", signature) // Required signature param
	if err != nil {
		return res, err
	}
	err = bodyWriter.WriteField("api_key", b.ApiKey) // Required signature param
	if err != nil {
		return res, err
	}

	err = bodyWriter.Close() // Close body writer
	if err != nil {
		return res, err
	}

	req, err := http.NewRequest("POST", url, bodyBuf) // Create new POST  request
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req) // Send request
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body) // Read response
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodyBytes, &res) // Unmarshal response in to ploadResponseData
	if err != nil {
		return res, err
	}
	if len(res.Error.Message) > 0 {
		return res, CustomError{res.Error.Message}
	}
	return res, nil
}

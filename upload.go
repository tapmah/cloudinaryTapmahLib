package cloudinarylib

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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
}

func (b *BaseCloudinaryStuct) UploadFile(filePath string, publicName string) (UploadResponseData, error) {

	var res UploadResponseData

	file, err := os.Open(filePath)
	if err != nil {
		return res, err
	}
	defer file.Close()

	// Создание параметров загрузки
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	params := map[string]string{
		"timestamp": timestamp,
		//"public_id": publicName,
	}
	sortedParams := sortParams(params) + b.ApiSecret

	sha := sha1.New()
	sha.Write([]byte(sortedParams))
	hash := sha.Sum(nil)
	sha1Str := fmt.Sprintf("%x", hash)

	url := "https://api.cloudinary.com/v1_1/" + b.CloudName + "/image/upload"

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", filePath)
	if err != nil {
		return res, err
	}
	file, err = os.Open(filePath)
	if err != nil {
		return res, err
	}
	defer file.Close()
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return res, err
	}

	err = bodyWriter.WriteField("api_key", b.ApiKey)
	if err != nil {
		return res, err
	}

	err = bodyWriter.WriteField("public_id", publicName)
	if err != nil {
		return res, err
	}
	err = bodyWriter.WriteField("timestamp", timestamp)
	if err != nil {
		return res, err
	}
	err = bodyWriter.WriteField("signature", sha1Str)
	if err != nil {
		return res, err
	}
	err = bodyWriter.Close()
	if err != nil {
		return res, err
	}

	req, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	// Чтение ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return res, err
	}

	// Вывод ответа
	return res, nil
}

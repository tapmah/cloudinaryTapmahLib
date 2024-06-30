package cloudinarylib

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
)

type BaseCloudinaryStuct struct {
	CloudName      string
	ApiKey         string
	ApiSecret      string
	Base64Key      string
	PublicID       string `json:"public_id"`
	TimeStamp      string `json:"timestamp"`
	Transformation string `json:"transformation"`
	Folder         string `json:"folder"`
}

type CropType string

const (
	CropFill  CropType = "fill"
	CropScale CropType = "scale"
	CropThumb CropType = "thumb"
	CropPad   CropType = "pad"
	CropLimit CropType = "limit"
)

type CustomError struct {
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}

func (b BaseCloudinaryStuct) Initialize(cloudName string, apiKey string, apiSecret string, publicID string) BaseCloudinaryStuct {
	return BaseCloudinaryStuct{
		CloudName: cloudName,
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		PublicID:  publicID,
		Base64Key: func() string {
			return base64.StdEncoding.EncodeToString([]byte(apiKey + ":" + apiSecret))
		}(),
	}
}

func (b *BaseCloudinaryStuct) SetCropParams(height string, width string, cropType CropType) {
	b.Transformation = fmt.Sprintf("c_%s,w_%s,h_%s", cropType, width, height)
}
func (b *BaseCloudinaryStuct) SetFolder(folder string) {
	b.Folder = folder
}

func sortParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	sortedParams := ""
	for _, key := range keys {
		sortedParams += key + "=" + params[key] + "&"
	}
	return sortedParams[:len(sortedParams)-1]
}

func ComputeSHA1(str string) string {
	sha := sha1.New()
	sha.Write([]byte(str))
	hash := sha.Sum(nil)
	sha1Str := fmt.Sprintf("%x", hash)
	return sha1Str
}

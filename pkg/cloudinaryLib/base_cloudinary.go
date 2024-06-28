package cloudinarylib

import "sort"

type BaseCloudinaryStuct struct {
	CloudName string
	ApiKey    string
	ApiSecret string
}

func (b BaseCloudinaryStuct) Initialize(cloudName string, apiKey string, apiSecret string) BaseCloudinaryStuct {
	return BaseCloudinaryStuct{
		CloudName: cloudName,
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
	}
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

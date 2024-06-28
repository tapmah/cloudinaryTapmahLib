package main

import (
	"log"

	cloudinarylib "github.com/tapmah/cloudinaryTapmahLib/pkg/cloudinaryLib"
)

func main() {

	const (
		apiKey    = "222882443985974"
		apiSecret = "3iGpikmfSByXW4hNN8IK1wsr3b4"
		cloudName = "dh8xbihlm"
	)

	cl := cloudinarylib.BaseCloudinaryStuct{}.Initialize(cloudName, apiKey, apiSecret)
	res, err := cl.UploadFile("D:/photo_2024-06-16_19-36-53.jpg", "KistiPhoto")
	if err != nil {
		log.Panicln(err)
	}
	print(res.SecureURL)
}

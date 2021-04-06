// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	rTypes "github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

var (
	bucketName string
	image      string
	client     *rekognition.Client
	checks     bool = true
)

func init() {

	flag.StringVar(&bucketName, "b", "", "The name of the bucket to get the object")
	flag.StringVar(&image, "i", "", "The path to the image file (JPEG, JPG, PNG)")
	flag.Parse()

	if len(bucketName) == 0 || len(image) == 0 {
		checks = false
		flag.PrintDefaults()
		log.Fatalf("You must supply a bucket name (-b BUCKET) and photo file (-i IMAGE)")
		return
	}

	fileExtension := filepath.Ext(image)
	validExtension := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
	}

	if !validExtension[fileExtension] {
		checks = false
		fmt.Println("Rekognition only supports jpeg, jpg or png")
		return
	}

	// Load the SDK's configuration from environment and shared config, and
	// create the client with this.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	client = rekognition.NewFromConfig(cfg)

}

func main() {

	if checks == false {
		return
	}

	params := &rekognition.DetectFacesInput{
		Image: &rTypes.Image{

			S3Object: &rTypes.S3Object{
				Bucket: &bucketName,
				Name:   &image,
			},
		},
		Attributes: []rTypes.Attribute{
			"ALL",
		},
	}

	resp, err := client.DetectFaces(context.TODO(), params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(resp.FaceDetails) == 0 {
		fmt.Println("No faces detected in the image !")
		return
	}

	for idx, fdetails := range resp.FaceDetails {

		fmt.Printf("Person #%d : \n", idx+1)
		fmt.Printf("Position : %v %v \n", *fdetails.BoundingBox.Left, *fdetails.BoundingBox.Top)

		if fdetails.AgeRange != nil {
			fmt.Printf("Age (Low) : %d \n", *fdetails.AgeRange.Low)
			fmt.Printf("Age (High) : %d \n", *fdetails.AgeRange.High)
		}

		if fdetails.Emotions != nil {
			fmt.Printf("Emotion : %v\n", fdetails.Emotions[0].Type)
		}

		if fdetails.Gender != nil {
			fmt.Printf("Gender : %v\n\n", fdetails.Gender.Value)
		}
	}

}
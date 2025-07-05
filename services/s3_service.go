package services

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
)

func ResizeImageFromS3(bucket, region string, key string, width string, height string) (string, error) {
	if width == "" || height == "" {
		return "", fmt.Errorf("width and height are required")
	}

	w, err := strconv.Atoi(width)
	if err != nil {
		return "", err
	}
	h, err := strconv.Atoi(height)
	if err != nil {
		return "", err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}
	defer output.Body.Close()

	img, format, err := image.Decode(output.Body)
	if err != nil {
		return "", err
	}

	resized := imaging.Resize(img, w, h, imaging.Lanczos)

	var buf bytes.Buffer
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, resized, nil)
	case "png":
		err = png.Encode(&buf, resized)
	default:
		return "", fmt.Errorf("unsupported image format: %s", format)
	}
	if err != nil {
		return "", err
	}

	newKey := "resized-" + key

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(newKey),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String("image/" + format),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, newKey)
	return url, nil
}

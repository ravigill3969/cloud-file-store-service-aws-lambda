package handlers

import (
	"context"
	"encoding/json"
	"go-lambda/services"

	"github.com/aws/aws-lambda-go/events"
)

func ImageHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	action := req.QueryStringParameters["action"]

	switch action {
	case "edit":
		url, err := services.ResizeImageFromS3(
			req.QueryStringParameters["bucketName"],
			req.QueryStringParameters["region"],
			req.QueryStringParameters["key"],
			req.QueryStringParameters["width"],
			req.QueryStringParameters["height"],
		)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}

		body, _ := json.Marshal(map[string]string{
			"message":   "Image processed",
			"image_url": url,
		})

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(body),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil

	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Action not supported",
		}, nil
	}
}

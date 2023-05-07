package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var table string
var client *dynamodb.Client

func init() {
	table = os.Getenv("TABLE_NAME")
	if table == "" {
		log.Fatal("missing environment variable TABLE_NAME")
	}
	cfg, _ := config.LoadDefaultConfig(context.Background())
	client = dynamodb.NewFromConfig(cfg)

}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	fmt.Println("function invoked with payload", req.Body)

	var p Payload
	err := json.Unmarshal([]byte(req.Body), &p)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	item := make(map[string]types.AttributeValue)

	item["email"] = &types.AttributeValueMemberS{Value: p.Email}
	item["name"] = &types.AttributeValueMemberS{Value: p.Name}

	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      item})

	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	log.Println("saved item with email", p.Email)

	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusOK}, nil
}

type Payload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

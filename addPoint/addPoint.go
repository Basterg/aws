package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

type Point struct {
	ExcursionID string `json:"excursionID"`
	Description string `json:"description"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	pointID, _ := uuid.NewV1()

	body := request.Body
	var thisPoint Point
	json.Unmarshal([]byte(body), &thisPoint)

	var pointToAdd = map[string]string{
		"PK":          fmt.Sprintf("EXC#%s", thisPoint.ExcursionID),
		"SK":          fmt.Sprintf("POINT#%s", pointID),
		"PointID":     fmt.Sprintf("%s", pointID),
		"Description": thisPoint.Description,
	}

	pointItem, err := dynamodbattribute.MarshalMap(pointToAdd)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Marshalling error", StatusCode: 500}, nil
	}

	_, err = ddb.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				ConditionCheck: &dynamodb.ConditionCheck{
					Key: map[string]*dynamodb.AttributeValue{
						"PK": {
							S: aws.String(fmt.Sprintf("EXC#%s", thisPoint.ExcursionID)),
						},
						"SK": {
							S: aws.String(fmt.Sprintf("EXC#%s", thisPoint.ExcursionID)),
						},
					},
					ConditionExpression: aws.String("attribute_exists(PK)"),
					TableName:           aws.String("Excursions"),
				},
			},
			{
				Put: &dynamodb.Put{
					TableName:           aws.String("Excursions"),
					Item:                pointItem,
					ConditionExpression: aws.String("attribute_not_exists(PK)"),
				},
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Error: %s", err), StatusCode: 200}, nil
	}
	return events.APIGatewayProxyResponse{Body: "Point added successfully", StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}

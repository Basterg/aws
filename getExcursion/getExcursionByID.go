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
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Excursion struct {
	ExcursionID string
	Country     string
	Description string
	Town        string
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	var thisExcursion Excursion

	var excID string
	excID = request.QueryStringParameters["excID"]

	proJ := expression.NamesList(
		expression.Name("ExcursionID"), expression.Name("Country"),
		expression.Name("Description"), expression.Name("Town"))

	expr, err := expression.NewBuilder().WithProjection(proJ).Build()
	if err != nil {
		fmt.Println(err)
	}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(fmt.Sprintf("EXC#%s", excID)),
			},
			"SK": {
				S: aws.String(fmt.Sprintf("EXC#%s", excID)),
			},
		},
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
		TableName:                aws.String("Excursions"),
	}

	result, err := ddb.GetItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Get - error", StatusCode: 200}, nil
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{Body: "Can't get the excursion", StatusCode: 200}, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &thisExcursion)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Unmarshalling error", StatusCode: 200}, nil
	}

	var excursionToSend = map[string]string{
		"excursionID": thisExcursion.ExcursionID,
		"country":     thisExcursion.Country,
		"description": thisExcursion.Description,
		"town":        thisExcursion.Town,
	}

	rez, _ := json.Marshal(excursionToSend)

	return events.APIGatewayProxyResponse{Body: string(rez), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}

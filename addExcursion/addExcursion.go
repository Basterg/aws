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
	"strings"
)

type Excursion struct {
	Country     string `json:"country"`
	Description string `json:"description"`
	Town        string `json:"town"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	body := request.Body

	excID, _ := uuid.NewV1()

	var thisExcursion Excursion
	json.Unmarshal([]byte(body), &thisExcursion)

	var excursionToAdd = map[string]string{
		"PK":          fmt.Sprintf("EXC#%s", excID),
		"SK":          fmt.Sprintf("EXC#%s", excID),
		"GSI1PK":      strings.ToUpper(thisExcursion.Country),
		"GSI1SK":      fmt.Sprintf("%s#EXC#%s", strings.ToUpper(thisExcursion.Town), excID),
		"Country":     thisExcursion.Country,
		"Description": thisExcursion.Description,
		"Town":        thisExcursion.Town,
		"ExcursionID": fmt.Sprintf("%s", excID),
	}

	av, err := dynamodbattribute.MarshalMap(excursionToAdd)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Marshalling error", StatusCode: 500}, nil
	}

	input := &dynamodb.PutItemInput{
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
		TableName:           aws.String("Excursions"),
	}

	_, err = ddb.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Put error", StatusCode: 200}, nil
	}
	return events.APIGatewayProxyResponse{Body: "Excursion added", StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}

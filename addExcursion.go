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
)


type Excursion struct {
	ExcursionID  string `json:"excursionID"`
	CountryUpper string `json:"countryUpper"`
	Country		 string `json:"country"`
	Description  string `json:"description"`
	TownUpper	 string `json:"townUpper"`
	Town 		 string `json:"town"`
}

func Handler(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	body := request.Body

	var thisExcursion Excursion
	json.Unmarshal([]byte(body), &thisExcursion)

	var excursionToAdd = map[string]string{
		"PK":  	 	   fmt.Sprintf("EXC#ID%s", thisExcursion.ExcursionID),
		"SK": 	 	   fmt.Sprintf("EXC#ID%s", thisExcursion.ExcursionID),
		"GSI1PK":  	   thisExcursion.CountryUpper,
		"GSI1SK":  	   fmt.Sprintf("%s#EXC#ID%s",thisExcursion.TownUpper, thisExcursion.ExcursionID),
		"Country": 	   thisExcursion.Country,
		"Description": thisExcursion.Description,
		"Town": 	   thisExcursion.Town,
	}

	av, err := dynamodbattribute.MarshalMap(excursionToAdd)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Post error", StatusCode: 500}, nil
	}

	input := &dynamodb.PutItemInput{
		Item:      			 av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
		TableName: 			 aws.String("Excursions"),
	}

	_, err = ddb.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Put error" , StatusCode: 200}, nil
	}
	return events.APIGatewayProxyResponse{Body: "Excursion added" , StatusCode: 200}, nil

}

func main()  {
	lambda.Start(Handler)
}




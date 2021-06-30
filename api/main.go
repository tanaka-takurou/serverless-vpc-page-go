package main

import (
	"os"
	"fmt"
	"log"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type APIResponse struct {
	Message string `json:"message"`
}

type VpcData struct {
	VpcId     string `json:"VpcId"`
	CidrBlock string `json:"CidrBlock"`
}

type SubnetData struct {
	SubnetId         string `json:"SubnetId"`
	CidrBlock        string `json:"CidrBlock"`
	AvailabilityZone string `json:"AvailabilityZone"`
}

type Response events.APIGatewayProxyResponse

const layout  string = "2006-01-02 15:04"
var ec2Client *ec2.Client

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "getvpcs" :
			res, e := getVpcs(ctx)
			if e != nil {
				err = e
			} else {
				jsonBytes, _ = json.Marshal(APIResponse{Message: res})
			}
		case "getsubnets" :
			res, e := getSubnets(ctx)
			if e != nil {
				err = e
			} else {
				jsonBytes, _ = json.Marshal(APIResponse{Message: res})
			}
		}
	}
	log.Print(request.RequestContext.Identity.SourceIP)
	if err != nil {
		log.Print(err)
		jsonBytes, _ = json.Marshal(APIResponse{Message: fmt.Sprint(err)})
		return Response{
			StatusCode: 500,
			Body: string(jsonBytes),
		}, nil
	}
	return Response {
		StatusCode: 200,
		Body: string(jsonBytes),
	}, nil
}

func getVpcs(ctx context.Context)(string, error) {
	if ec2Client == nil {
		ec2Client = ec2.NewFromConfig(getConfig(ctx))
	}
	input := &ec2.DescribeVpcsInput{
		VpcIds: []string{os.Getenv("VPC_ID")},
	}

	result, err := ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		log.Print(err)
		return "", err
	}
	if len(result.Vpcs) < 1 {
		err := fmt.Errorf("No Vpcs")
		log.Print(err)
		return "", err
	}
	resultJson, err := json.Marshal(VpcData{
		VpcId: aws.ToString(result.Vpcs[0].VpcId),
		CidrBlock: aws.ToString(result.Vpcs[0].CidrBlock),
	})
	if err != nil {
		log.Print(err)
		return "", err
	}
	return string(resultJson), nil
}

func getSubnets(ctx context.Context)(string, error) {
	if ec2Client == nil {
		ec2Client = ec2.NewFromConfig(getConfig(ctx))
	}
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []string{os.Getenv("SUBNET_ID")},
	}

	result, err := ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		log.Print(err)
		return "", err
	}
	if len(result.Subnets) < 1 {
		err := fmt.Errorf("No Subnets")
		log.Print(err)
		return "", err
	}
	resultJson, err := json.Marshal(SubnetData{
		SubnetId: aws.ToString(result.Subnets[0].SubnetId),
		CidrBlock: aws.ToString(result.Subnets[0].CidrBlock),
		AvailabilityZone: aws.ToString(result.Subnets[0].AvailabilityZone),
	})
	if err != nil {
		log.Print(err)
		return "", err
	}
	return string(resultJson), nil
}

func getConfig(ctx context.Context) aws.Config {
	var err error
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Print(err)
	}
	return cfg
}

func main() {
	lambda.Start(HandleRequest)
}

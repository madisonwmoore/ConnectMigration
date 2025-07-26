package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type ContactFlowOptions struct {
	name        string
	description string
	instance_id string
	region      *string
	flowType    string
	content     string
	tags        *[]string
}

func getContactFlow(flowId string, instanceId string) *connect.DescribeContactFlowOutput {
	sdkConfig, _ := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("ccc-dev"))
	client := connect.NewFromConfig(sdkConfig)
	var params connect.DescribeContactFlowInput
	params.ContactFlowId = &flowId
	params.InstanceId = &instanceId
	var output *connect.DescribeContactFlowOutput
	output, err := client.DescribeContactFlow(context.TODO(), &params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*output.ContactFlow.Content)

	return output
}

func addContactFlow(name string, rootBody *hclwrite.Body, config ContactFlowOptions) {
	fmt.Println("Adding Routing Profile")
	rp := rootBody.AppendNewBlock("resource", []string{"aws_connect_contact_flow", name})
	rp.Body().SetAttributeValue("name", cty.StringVal(config.name))
	rp.Body().SetAttributeValue("description", cty.StringVal(config.description))
	rp.Body().SetAttributeValue("instance_id", cty.StringVal(config.instance_id))
	if config.region != nil {
		rp.Body().SetAttributeValue("region", cty.StringVal(*config.region))
	}
	rp.Body().SetAttributeValue("type", cty.StringVal(config.flowType))
	rp.Body().SetAttributeRaw("content", hclwrite.TokensForFunctionCall("templatefile", hclwrite.TokensForFunctionCall("file")))
	//rp.Body().SetAttributeValue("content", cty.StringVal(config.content))
	//rp.Body().SetAttributeValue("tags", cty.StringVal(config.tags))
	rootBody.AppendNewline()
}

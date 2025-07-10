package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func main() {

	//getRoutingProfiles("53ffe28d-757b-4ca0-99d8-4babf0e1fe0f")

	createProviderFile("us-west-2")

	f := hclwrite.NewEmptyFile()
	tfFile, err := os.Create("routing_profiles.tf")
	if err != nil {
		fmt.Println(err)
		return
	}
	routingProfileBody := f.Body()
	addRoutingProfile("agents", routingProfileBody, RoutingProfileOptionsStruct{name: "agents", description: "agents description"})
	region := "us-east-1"
	addRoutingProfile("agents", routingProfileBody, RoutingProfileOptionsStruct{name: "agents", description: "agents description", region: &region})

	// for i := 0; i < 10; i++ {
	// 	addRoutingProfile(strconv.Itoa(i), rootBody)
	// }

	tfFile.Write(f.Bytes())

}

type media_concurrencies struct {
	channel     string
	concurrency int
}

type RoutingProfileOptionsStruct struct {
	name          string
	description   string
	region        *string
	media_channel *[]media_concurrencies
}

func createProviderFile(region string) *hclwrite.Body {
	f := hclwrite.NewEmptyFile()
	tfFile, err := os.Create("providers.tf")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	providersBody := f.Body()
	terraformBlock := providersBody.AppendNewBlock("terraform", nil)
	requiredProvidersBlock := terraformBlock.Body().AppendNewBlock(" required_providers", nil)
	requiredProvidersBlock.Body().SetAttributeValue("aws", cty.ObjectVal(map[string]cty.Value{
		"source":  cty.StringVal("hashicorp/aws"),
		"version": cty.StringVal("~> 6.0"),
	}))

	providerBlock := providersBody.AppendNewBlock("provider", []string{"aws"})
	providerBlock.Body().SetAttributeValue("region", cty.StringVal(region))
	tfFile.Write(f.Bytes())
	return providersBody
}

func RoutingProfileOptions(name string, description string) RoutingProfileOptionsStruct {
	rpo := RoutingProfileOptionsStruct{}
	rpo.name = name
	rpo.description = description
	return rpo
}

func getRoutingProfiles(instance_id string) {
	// ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("collectors-staging"))
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	client := connect.NewFromConfig(sdkConfig)

	InstanceId := instance_id
	result, err := client.ListRoutingProfiles(context.TODO(), &connect.ListRoutingProfilesInput{InstanceId: &InstanceId})

	if err != nil {
		fmt.Println(err)
	}

	for _, profile := range result.RoutingProfileSummaryList {
		fmt.Printf("  - Name: %s, ID: %s\n", *profile.Name, *profile.Id)
	}

}

func addRoutingProfile(name string, rootBody *hclwrite.Body, config RoutingProfileOptionsStruct) {
	rp := rootBody.AppendNewBlock("resource", []string{"aws_connect_routing_profile", name})
	rp.Body().SetAttributeValue("name", cty.StringVal("w"))
	rp.Body().SetAttributeValue("description", cty.StringVal("w"))
	if config.region != nil {
		rp.Body().SetAttributeValue("region", cty.StringVal(*config.region))
	}

	rp.Body().SetAttributeValue("default_outbound_queue_id", cty.StringVal("world"))
	rp.Body().SetAttributeValue("instance_id", cty.StringVal("world"))

	mc := rp.Body().AppendNewBlock("media_concurrencies", nil)
	mc.Body().SetAttributeValue("channel", cty.StringVal("VOICE"))
	mc.Body().SetAttributeValue("concurrency", cty.NumberIntVal(1))

	rootBody.AppendNewline()
}

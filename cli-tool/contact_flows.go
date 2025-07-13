package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type ContactFlowOptions struct {
	name        string
	description string
	region      *string
	flowType    string
	content     string
	tags        *[]string
}

func addContactFlow(name string, rootBody *hclwrite.Body, config ContactFlowOptions) {
	fmt.Println("Adding Routing Profile")
	rp := rootBody.AppendNewBlock("resource", []string{"aws_connect_routing_profile", name})
	rp.Body().SetAttributeValue("name", cty.StringVal(config.name))
	rp.Body().SetAttributeValue("description", cty.StringVal(config.description))
	if config.region != nil {
		rp.Body().SetAttributeValue("region", cty.StringVal(*config.region))
	}
	rp.Body().SetAttributeValue("type", cty.StringVal(config.flowType))
	rp.Body().SetAttributeValue("content", cty.StringVal(config.content))
	rootBody.AppendNewline()
}

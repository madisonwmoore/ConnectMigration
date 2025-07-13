package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type RoutingProfileOptionsStruct struct {
	name          string
	description   string
	region        *string
	media_channel *[]media_concurrencies
}

type media_concurrencies struct {
	channel     string
	concurrency int
}

func addRoutingProfile(name string, rootBody *hclwrite.Body, config RoutingProfileOptionsStruct) {
	fmt.Println("Adding Routing Profile")
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

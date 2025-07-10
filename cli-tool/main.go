package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func main() {
	f := hclwrite.NewEmptyFile()
	tfFile, err := os.Create("contact_flow.tf")
	if err != nil {
		fmt.Println(err)
		return
	}
	rootBody := f.Body()
	addRoutingProfile("agents", rootBody, RoutingProfileOptions("agents", "agents description"))
	// addRoutingProfile("supervisor", rootBody)

	// for i := 0; i < 10; i++ {
	// 	addRoutingProfile(strconv.Itoa(i), rootBody)
	// }

	tfFile.Write(f.Bytes())

}

type RoutingProfileOptionsStruct struct {
	name        string
	description string
	region      string
}

func RoutingProfileOptions(name string, description string) RoutingProfileOptionsStruct {
	rpo := RoutingProfileOptionsStruct{}
	rpo.name = name
	rpo.description = description
	return rpo
}

func addRoutingProfile(name string, rootBody *hclwrite.Body, config *RoutingProfileOptionsStruct) {
	rp := rootBody.AppendNewBlock("resource", []string{"aws_connect_routing_profile", name})
	rp.Body().SetAttributeValue("name", cty.StringVal("w"))
	rp.Body().SetAttributeValue("description", cty.StringVal("w"))
	if config.region != nil {
		rp.Body().SetAttributeValue("region", cty.StringVal("w"))
	}

	rp.Body().SetAttributeValue("default_outbound_queue_id", cty.StringVal("world"))
	rp.Body().SetAttributeValue("instance_id", cty.StringVal("world"))

	mc := rp.Body().AppendNewBlock("media_concurrencies", nil)
	mc.Body().SetAttributeValue("channel", cty.StringVal("VOICE"))
	mc.Body().SetAttributeValue("concurrency", cty.NumberIntVal(1))

	rootBody.AppendNewline()
}

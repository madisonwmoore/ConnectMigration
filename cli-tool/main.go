package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type LexBots struct {
	profile   string
	ArnToName map[string]string
	NameToArn map[string]string
}

func (l LexBots) set(name string, arn string) {
	l.ArnToName[arn] = name
	l.NameToArn[name] = arn
}

func (l LexBots) getARN(name string) string {
	return l.NameToArn[name]
}

func (l LexBots) getName(arn string) string {
	return l.ArnToName[arn]
}

func main() {
	//lexBotStore := LexBots{profile: "ccc-dev", ArnToName: make(map[string]string), NameToArn: make(map[string]string)}
	//getLexBots("ccc-dev", &lexBotStore)

	queue_map := make(map[string]map[string]string)
	// getQueues("ccc-dev", "d5826fea-981d-4b94-8371-40dff206ea5b", queue_map)
	parseContactFlow(queue_map, "ccc-dev")

	// uat_queues := getQueues("ccc-uat", "85c504f2-ea80-4970-bc37-7073b2b8dacb", queue_map)
	// for key, value := range uat_queues {
	// 	fmt.Printf("Key: %s, Value: %s\n", key, value)
	// }

	//parseContactFlow()

	instanceId := "d5826fea-981d-4b94-8371-40dff206ea5b"
	output := getContactFlow("2f4eb27d-550b-4d12-aa9b-f31d471d2148", instanceId)

	createProviderFile("us-west-2")

	f := hclwrite.NewEmptyFile()
	tfFile, err := os.Create("out/routing_profiles.tf")
	if err != nil {
		fmt.Println(err)
		return
	}
	routingProfileBody := f.Body()
	// addRoutingProfile("agents", routingProfileBody, RoutingProfileOptionsStruct{name: "agents", description: "agents description"})
	region := "us-east-1"
	// addRoutingProfile("agents", routingProfileBody, RoutingProfileOptionsStruct{name: "agents", description: "agents description", region: &region})

	// addContactFlow("boo", routingProfileBody, nil)

	addContactFlow(*output.ContactFlow.Arn, routingProfileBody, ContactFlowOptions{content: *output.ContactFlow.Content, name: "New_Flow", region: &region})

	// for i := 0; i < 10; i++ {
	// 	addRoutingProfile(strconv.Itoa(i), rootBody)
	// }

	tfFile.Write(f.Bytes())

}

func createProviderFile(region string) *hclwrite.Body {
	f := hclwrite.NewEmptyFile()
	tfFile, err := os.Create("out/providers.tf")
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

func parseContactFlow(queue_map map[string]map[string]string, profile string) {
	fmt.Println("Parsing Contact Flow")
	m := make(map[string]string)

	type Parameters struct {
		QueueId  string `json:"QueueId"`
		LexV2Bot struct {
			AliasArn string `json:"AliasArn"`
		} `json:"LexV2Bot"`
	}

	type Action struct {
		Identifier string     `json:"Identifier"`
		Type       string     `json:"Type"`
		Parameters Parameters `json:"Parameters"`
	}
	type Content struct {
		Actions  []Action `json:"Actions"`
		Metadata struct{} `json:"Metadata"`
	}

	jsonData, err := os.ReadFile("./LEXINTENT_ACCOUNT_BALANCE.json.tpl")
	if err != nil {
		fmt.Print(err)
		// Handle error
	}
	var c Content
	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		fmt.Println(err)
		// Handle error
	}

	for i := 0; i < len(c.Actions); i++ {
		var action Action = c.Actions[i]
		if action.Type == "UpdateContactTargetQueue" {
			queueName := "UNKNOWN"
			for queueNameKey, value := range queue_map {
				if value[profile] == c.Actions[i].Parameters.QueueId {
					queueName = queueNameKey
					break
				}
			}
			m[c.Actions[i].Parameters.QueueId] = "QUEUE_ARN_" + queueName
		}
		if action.Type == "ConnectParticipantWithLexBot" {
			lexAliasName := "UNKNOWN"
			m[c.Actions[i].Parameters.LexV2Bot.AliasArn] = "LEX_ALIAS_ARN_" + lexAliasName
			//fmt.Println(c.Actions[i].Parameters.LexV2Bot.AliasArn)
		}
	}

	for key, value := range m {
		jsonData = []byte(strings.ReplaceAll(string(jsonData), key, fmt.Sprintf("${%s}", value)))
	}

	err = os.WriteFile("output.json", []byte(jsonData), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getQueues(profile string, instance_id string, queue_map map[string]map[string]string) map[string]map[string]string {
	// queue_map := make(map[string]interface{})
	instanceID := instance_id
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}
	client := connect.NewFromConfig(cfg)

	// Prepare the input for ListQueues
	input := &connect.ListQueuesInput{
		InstanceId: &instanceID,
		// You can filter by QueueTypes if needed, e.g., []types.QueueType{types.QueueTypeStandard}
		// QueueTypes: []types.QueueType{types.QueueTypeStandard},
	}

	// Call the ListQueues API
	resp, err := client.ListQueues(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to list queues, %v", err)
	}

	for _, queue := range resp.QueueSummaryList {
		if queue.QueueType == "STANDARD" {
			var queueName string = strings.ReplaceAll(strings.ToUpper(*queue.Name), " ", "_")
			if queue_map[*queue.Name] == nil {
				queue_map[queueName] = make(map[string]string)
			}
			queue_map[queueName][profile] = *queue.Arn
		}
	}
	//fmt.Printf("ARN: %s, Type: %s\n", *queue.Arn, queue.QueueType)
	return queue_map
}

func getLexBots(profile string, store *LexBots) {
	fmt.Println(store.profile)
	// lexBotMap := make(map[string]map[string]string)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}
	client := lexmodelsv2.NewFromConfig(cfg)

	input := &lexmodelsv2.ListBotsInput{}

	// Call the ListBots operation
	resp, err := client.ListBots(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to list bots, %v", err)
	}

	fmt.Println("Amazon Lex Bots:")
	for _, bot := range resp.BotSummaries {
		fmt.Printf("  Bot Name: %s, Bot ID: %s\n", *bot.BotName, *bot.BotId)
		// store.[*bot.BotName] = make(map[string]string)
	}

	// Handle pagination if more results are available
	for resp.NextToken != nil {
		input.NextToken = resp.NextToken
		resp, err = client.ListBots(context.TODO(), input)
		if err != nil {
			log.Fatalf("failed to list bots (pagination), %v", err)
			break
		}
		for _, bot := range resp.BotSummaries {
			input := &lexmodelsv2.ListBotAliasesInput{}
			input.BotId = bot.BotId
			resp, err := client.ListBotAliases(context.TODO(), input)
			if err != nil {
				log.Fatalf("failed to list bots (pagination), %v", err)
			}
			for _, alias := range resp.BotAliasSummaries {
				fmt.Println(alias.BotAliasId)

			}
			for resp.NextToken != nil {
				input.NextToken = resp.NextToken
				resp, err := client.ListBotAliases(context.TODO(), input)
				if err != nil {
					log.Fatalf("failed to list bots (pagination), %v", err)
				}
				for _, alias := range resp.BotAliasSummaries {
					store.set(fmt.Sprintf("%s/%s", *bot.BotName, *alias.BotAliasName), fmt.Sprintf("arn:aws:lex:%s:%s:bot-alias/%s/%s", *bot.BotName, *alias.BotAliasName))
					fmt.Println(alias.BotAliasId)
				}
			}
			fmt.Println(resp.BotAliasSummaries)
			fmt.Printf("Bot Name: %s, Bot ID: %s\n", *bot.BotName, *bot.BotId)
		}
	}

	//return lexBotMap
}

// func RoutingProfileOptions(name string, description string) RoutingProfileOptionsStruct {
// 	rpo := RoutingProfileOptionsStruct{}
// 	rpo.name = name
// 	rpo.description = description
// 	return rpo
// }

// func getRoutingProfiles(instance_id string) {
// 	// ctx := context.Background()
// 	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("collectors-staging"))
// 	if err != nil {
// 		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
// 		fmt.Println(err)
// 		return
// 	}
// 	client := connect.NewFromConfig(sdkConfig)

// 	InstanceId := instance_id
// 	result, err := client.ListRoutingProfiles(context.TODO(), &connect.ListRoutingProfilesInput{InstanceId: &InstanceId})

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	for _, profile := range result.RoutingProfileSummaryList {
// 		fmt.Printf("  - Name: %s, ID: %s\n", *profile.Name, *profile.Id)
// 	}

// }

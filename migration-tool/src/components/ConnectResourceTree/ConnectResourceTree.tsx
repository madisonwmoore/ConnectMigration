import {
  Box,
  Tree,
  Text,
  Badge,
  Group,
  ActionIcon,
  Stack,
  Paper,
  ScrollArea,
  Checkbox,
} from "@mantine/core";
import {
  IconArrowLeft,
  IconBinaryTree,
  IconLambda,
  IconRefresh,
  IconRobot,
} from "@tabler/icons-react";
import { Children } from "react";
import type { TreeNodeData } from '@mantine/core';

// CONTACT_FLOW | CUSTOMER_QUEUE | CUSTOMER_HOLD | CUSTOMER_WHISPER | AGENT_HOLD | AGENT_WHISPER | OUTBOUND_WHISPER | AGENT_TRANSFER | QUEUE_TRANSFER | CAMPAIGN          


function ConnectResourceTree(props:any) {
  const createContactFlowResources:(name: string)=>TreeNodeData = (name: string) => {
    return {
      label: <Group><Checkbox.Indicator/><Text c="indigo">{name}</Text></Group>,
      value: "Contact Flow Name",
      children: [
        {
          value: "src/components",
          label: (
            <span>
              <Text fw={500}>
                <IconRobot color="orange"/>
                Lex Bots
              </Text>
            </span>
          ),
          children: [],
        },
        {
          value: "src/hooks",
          label: (
            <span>
              <Text fw={"500"}>
                <IconLambda color="orange" />
                Lambda Functions
              </Text>
            </span>
          ),
          children: [{ value: "src/components", label: "Lambda Function" }],
        },
      ],
    } as TreeNodeData;
  };

  const data = [
    {
      value: "src",
      label: (
        <>
          <IconBinaryTree color="orange" />
          Contact Flows
        </>
      ),
      children: [{label:<Text size="xs" >Inbound Flow</Text>},createContactFlowResources("Contact Flow 1"),{label:<Text size="xs" >Disconnect Flow</Text>},createContactFlowResources("Contact Flow 1"),{label:<Text size="xs" >Agent Whisper Flow</Text>},createContactFlowResources("Contact Flow 1")],
    },
    { label: "Prompts", value: "prompts" },
  ];

  return (

    <Paper withBorder p={"md"} w={500}>
    <ScrollArea h={600}>
      <Group justify="space-between">
        <Badge>Instance Name</Badge>
        <ActionIcon variant="outline">
          <IconRefresh />
        </ActionIcon>
      </Group>
      <Tree ta={"left"} levelOffset={10} data={data} />
      </ScrollArea>
    </Paper>
  );
}

export default ConnectResourceTree;

resource "aws_connect_routing_profile" "agents" {
  name                      = "w"
  description               = "w"
  default_outbound_queue_id = "world"
  instance_id               = "world"
  media_concurrencies {
    channel     = "VOICE"
    concurrency = 1
  }
}

resource "aws_connect_routing_profile" "agents" {
  name                      = "w"
  description               = "w"
  region                    = "us-east-1"
  default_outbound_queue_id = "world"
  instance_id               = "world"
  media_concurrencies {
    channel     = "VOICE"
    concurrency = 1
  }
}

resource "aws_connect_contact_flow" "boo" {
  name = "w"
}

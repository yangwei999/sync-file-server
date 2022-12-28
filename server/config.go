package main

import "github.com/opensourceways/community-robot-lib/mq"

type configuration struct {
	MQConfig mq.MQConfig `json:"mq_config" required:"true"`
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/wire"
)

func main() {

	var di = bootstrap.Get()

	app, err := wire.InitializeCdcConsumer(di)
	if err != nil {
		panic(err)
	}

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	addr := fmt.Sprintf("%s:%s", di.Env.Kafka.Address, di.Env.Kafka.Port)

	group, err := sarama.NewConsumerGroup([]string{addr}, di.Const.Kafka.GroupidForElastic, config)

	if err != nil {
		panic(err)
	}
	defer group.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		err := group.Consume(ctx, di.Const.Kafka.Cdctopics, app.KafkaCdc)
		if err != nil {
			log.Fatalf("Consumer error: %v", err)
		}
		if ctx.Err() != nil {
			break
		}
	}
}

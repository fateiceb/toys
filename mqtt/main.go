package main

import (
	mq "github.com/eclipse/paho.mqtt.golang"
	"time"
)
func main() {
	options := mq.NewClientOptions().SetClientID("1").AddBroker("tcp://222.24.22.223")
	options.SetKeepAlive(2 * time.Second)
	//options.SetDefaultPublishHandler(messageHandler)
	options.SetPingTimeout(1 * time.Second)
	client := mq.NewClient(options)
	token := client.Connect()
	if token.Wait() {
	}

	select {

	}
}

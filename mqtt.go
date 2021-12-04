package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"strconv"
	"time"
)

type MqttBus struct {
	Debug     bool
	Server    string
	Port      int
	Username  string
	Password  string
	ClientId  string
	TopicRoot string
	client    MQTT.Client
}

func (mqtt *MqttBus) Initialize() {
	if mqtt.Debug {
		fmt.Println("INITIALIZING MQTT BUS...")
	}
	if mqtt.Port == 0 {
		mqtt.Port = 1883
	}
	if mqtt.TopicRoot == "" {
		mqtt.TopicRoot = "bobcat-monitor"
	}
	mqttOpts := MQTT.NewClientOptions().AddBroker("tcp://" + mqtt.Server + ":" + strconv.Itoa(mqtt.Port))
	mqttOpts.SetUsername(mqtt.Username)
	if mqtt.Password != "" {
		mqttOpts.SetPassword(mqtt.Password)
	}
	mqttOpts.SetAutoReconnect(true)
	if mqtt.ClientId == "" {
		mqttOpts.SetClientID("bobcat-monitor-" + strconv.Itoa(rand.Intn(100)))
	} else {
		mqttOpts.SetClientID(mqtt.ClientId)
	}
	mqttOpts.SetKeepAlive(2 * time.Second)
	mqttOpts.SetPingTimeout(1 * time.Second)
	mqttOpts.SetWill(mqtt.TopicRoot+"/monitor", `{ "status": "down" }`, 0, false)

	mqttOpts.OnConnect = func(client MQTT.Client) {
		fmt.Printf("MQTT: CONNECTED TO %s\n", mqtt.Server)
	}

	mqttOpts.DefaultPublishHandler = func(client MQTT.Client, msg MQTT.Message) {
		if mqtt.Debug {
			fmt.Printf("  MQTT: TOPIC: %s\n  MQTT: MESSAGE: %s\n", msg.Topic(), msg.Payload())
		}
	}

	mqtt.client = MQTT.NewClient(mqttOpts)
	if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqtt.SendMessage(mqtt.TopicRoot+"/monitor", `{ "status": "up" }`)
}

func (mqtt *MqttBus) SendMessage(topic string, payload interface{}) {
	if !mqtt.client.IsConnected() {
		fmt.Println("MQTT: CLIENT NOT CONNECTED")
		return
	}
	if token := mqtt.client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		fmt.Printf("MQTT ERROR, %s\n", token.Error())
	}
}

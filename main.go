package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Debug           bool   `yaml:"debug"`
	BobcatAddress   string `yaml:"bobcatAddress"`
	IntervalSeconds int    `yaml:"intervalSeconds"`
	Mqtt            struct {
		Enabled   bool   `yaml:"enabled"`
		Server    string `yaml:"server"`
		Port      int    `yaml:"port,omitempty"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		ClientId  string `yaml:"clientId,omitempty"`
		TopicRoot string `yaml:"topicRoot"`
	} `yaml:"mqtt"`
}

func readConfig() (*Config, error) {
	config := &Config{}
	file, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, os.Interrupt, syscall.SIGTERM)

	// LOAD CONFIG
	config, err := readConfig()
	if err != nil {
		log.Fatalln("Failed to read config. Exiting!")
		return
	}
	mqttBus := MqttBus{}

	if config.Mqtt.Enabled {
		if config.Debug {
			fmt.Println("MQTT ENABLED, CONNECTING...")
		}
		if config.Mqtt.TopicRoot == "" {
			config.Mqtt.TopicRoot = "bobcat-monitor"
		}
		mqttBus.Debug = config.Debug
		mqttBus.Server = config.Mqtt.Server
		mqttBus.Port = config.Mqtt.Port
		mqttBus.Username = config.Mqtt.Username
		mqttBus.Password = config.Mqtt.Password
		mqttBus.ClientId = config.Mqtt.ClientId
		mqttBus.TopicRoot = config.Mqtt.TopicRoot
		mqttBus.Initialize()
		if config.Debug {
			fmt.Println("MQTT BUS INITIALIZED")
		}
	}

	bobcatChannel := make(chan BobcatStatus, 5)

	bobcat := Bobcat{
		Debug:         config.Debug,
		periodSeconds: config.IntervalSeconds,
		address:       config.BobcatAddress,
		eventChannel:  bobcatChannel,
	}

	go bobcat.Begin()

	println("BOBCAT INIT DONE!")

	go func() {
		for {
			bobcatStatus := <-bobcatChannel
			if config.Mqtt.Enabled {
				if config.Debug {
					log.Println("POSTING STATUS TO MQTT...")
				}
				payload, err := json.Marshal(bobcatStatus)
				if err != nil {
					log.Printf("Error while unmarshalling JSON: %s", err)
				} else {
					var topic = "bobcat"
					mqttBus.SendMessage(config.Mqtt.TopicRoot+"/"+topic, payload)
					if config.Debug {
						log.Printf("SENT MQTT MESSAGE: %s TO TOPIC %s \n", payload, topic)
					}
				}
			}
		}
	}()

	// WAIT FOR SIGTERM FOREVER
	<-osChannel
}

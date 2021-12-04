# Bobcat Monitor

Small app to keep tabs on your Bobcat Helium miner's health and integrate it with anything through MQTT.

## What does it do?

It scrapes the status of your miner and posts it to MQTT. Really simple!

It enables many scenarios like
  
  - Trigger a smart Wi-Fi outlet on and off when your miner is stuck offline (thanks, Bobcat!)
  - [Ingest data into InfluxDB](https://www.influxdata.com/integration/mqtt-monitoring/) to create pretty graphs
  - Send you email/discord/telegram messages when your miner goes too far out of sync

Note that none of these scenarios are in scope of this application and have to be implemented by other software components - like [Home Assistant](https://www.home-assistant.io) or [Node-RED](https://nodered.org).

## Running with Docker

You can either use ready-built image or build your own.

`docker run -d -v $PWD/config.yml:/config.yml quay.io/toxuin/bobcat-monitor`

Pre-built image supports x86_64, armv6, armv7, arm64v8, and ppc64le.

Works on:
- Normal computers
- Raspberry pi
- ARM macs
- Many more things!

## Configuration

Sample config (config.yml) contains all the possible configuration options, which are pretty self-explanatory.

`debug`: if set to true, will make app more verbose. Helps when troubleshooting issues. Default: false

`bobcatAddress`: that is ip or domain name (if you have one assigned) of your miner. Required.

`intervalSeconds`: how often should this app check on your bobcat? Note that sometimes (for example right after an OTA update) requests can take up to 10 seconds, so setting this too low is not recommended.

`mqtt`: section that has all the parameters of your MQTT broker

`enabled`: setting this to false will disable all the MQTT functionality

`server`: address of your MQTT broker. Required.

`port`: port of your MQTT broker. Will use 1883 if not set.

`username`: username to access your MQTT broker. For anonymous access, leave out of config file.

`password`: well... Password. For anonymous access - leave out of the file.

`clientId`: if your MQTT broker requires client id to be set - set it here. When not specified - will generate one.

`topicRoot`: topic under which all the bobcat-related messages will be posted. Defaults to `bobcat-monitor` when not set.


## Support the author

Please?

HNT: 14dwpLQQ6CkKFFuAnwbBQPUjkcnZeNK6o66zVLkvwvo8eybrUvx

BTC: bc1q84y9cxyeg2940unavvzx2j5r5lngge802jsrav

Thank you!
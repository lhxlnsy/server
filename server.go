package server

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	Redis "github.com/lhxlnsy/redis"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/panjf2000/ants/v2"
)

var messagecount int
var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("received message: %s\n", msg)
	messagecount++
	ConvertMessageToStruct(redis, msg.Payload(), msg.Topic())
}

type Mqtt interface {
	Publish(topic, payload string)
	Subscribe()
}

var redis *Redis.PAPRedis

type PAPMqtt struct {
	mqtt   Mqtt
	topic  string
	client mqtt.Client
}

func (m *PAPMqtt) Subcribe() {
	fmt.Printf("Start to subscribe topic: %v \n", m.topic)
	token := m.client.Subscribe(m.topic, 0, f)
	if token.Error() != nil {
		panic(token.Error())
	}
	token.Wait()
}

func (m *PAPMqtt) SetRedis(redispopl *Redis.PAPRedis) {
	redis = redispopl
}

func (m *PAPMqtt) Publish(args ...string) {

	var publishtopic string
	var payload string
	if len(args) == 1 {
		publishtopic = m.topic
		payload = args[0]
	} else {
		publishtopic = args[0]
		payload = args[1]
	}
	token := m.client.Publish(publishtopic, 0, false, payload)
	token.Wait()
	return
}

func (m *PAPMqtt) Close() {
	m.client.Disconnect(250)
}

func (m *PAPMqtt) Client() *mqtt.Client {
	return &m.client
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Printf("Connected\n")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	mqtt.ERROR.Printf("Connect lost: %v", err)
}

func NewMqtt(pool *ants.Pool, wg *sync.WaitGroup) *PAPMqtt {
	messagecount = 0
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://192.168.0.222:1883")
	opts.SetClientID("emqtt_planetarkpower_client_" + time.Now().String())
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.CleanSession = false
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return &PAPMqtt{
		topic:  "meter_grid/state",
		client: client,
	}
}

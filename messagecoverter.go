package server

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	models "github.com/lhxlnsy/server/models/meter_grid"

	Redis "github.com/lhxlnsy/redis"
)

var message map[string]interface{}

func ConvertMessageToStruct(redis *Redis.PAPRedis, msg []byte, topic string) {
	json.Unmarshal(msg, &message)
	for key, value := range message {
		//fmt.Printf("Reading Value %v for Key %v \n", value, key)
		topickey := topic + "_" + key
		redis.PushData(topickey, value)
	}
	// TODO: Check the redis list length here, if timepstamp len more than default length, then store to the database
	currentLen, err := redis.GetLen(topic + "_Timestamp")
	fmt.Printf("current len is :%d\n", currentLen)
	if err != nil {
		fmt.Printf("Error when get length %s\n", err)
	}
	if currentLen >= redis.GetDefaultLen() {
		fmt.Printf("current len is :%d\n", currentLen)
		fmt.Printf("Default len is :%d\n", redis.GetDefaultLen())
		StoreMessageToPostGres(topic)
	}
}

func StoreMessageToPostGres(topic string) {
	Data := make(map[string]interface{})
	for key, _ := range message {
		if key != "Timestamp" {
			storelist := redis.GetData(topic + "_" + key)
			redis.EmptyList(topic + "_" + key)
			var totalvalue float64
			totalvalue = 0
			for _, value := range storelist {
				v, _ := strconv.ParseFloat(value, 64)
				totalvalue = totalvalue + float64(v)
			}
			average := totalvalue / float64(redis.GetDefaultLen())
			Data[key] = math.Round(average*1000) / 1000
			fmt.Println("average:", average)
		} else {
			redis.EmptyList(topic + "_" + key)
			Data[key] = time.Now()
		}
	}
	fmt.Println(Data)
	jsonString, err := json.Marshal(Data)
	fmt.Println(jsonString)
	if err != nil {
		fmt.Println(err)
	}
	DataModel := &models.Meter_grid_stat{}
	json.Unmarshal(jsonString, &DataModel)
	fmt.Println(string(jsonString))
	fmt.Println("final Data:", DataModel)
	GormDb.Model(&models.Meter_grid_stat{}).Create(DataModel)
	fmt.Println("final Data:", DataModel)
}

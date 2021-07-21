module github.com/lhxlnsy/server

go 1.16

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/lhxlnsy/redis v0.0.1
	github.com/panjf2000/ants/v2 v2.4.6
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.9
	models v0.0.1
)

replace models v0.0.1 => ./models

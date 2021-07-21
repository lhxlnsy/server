package server

import (
	"fmt"

	models "github.com/lhxlnsy/server/models/meter_grid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgreOption struct {
	Host     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
	Port     string
}

func generatePostgreConnStr(option *postgreOption) string {
	return fmt.Sprintf(`host=%s user=%s password=%s port=%s timezone=%s sslmode=%s dbname=%s`,
		option.Host,
		option.User,
		option.Password,
		option.Port,
		option.TimeZone,
		option.SSLMode,
		option.DBName,
	)
}

var GormDb *gorm.DB

func Init() {
	dsn := &postgreOption{
		Host:     "192.168.0.222",
		User:     "planetarkpower",
		Password: "PAP2021",
		Port:     "5432",
		TimeZone: "Asia/Shanghai",
		SSLMode:  "disable",
		DBName:   "planetarkpower",
	}
	connstr := generatePostgreConnStr(dsn)
	db, err := gorm.Open(postgres.Open(connstr), &gorm.Config{})
	if err != nil {
		GormDb = nil
	}
	GormDb = db
	db.AutoMigrate(&models.Meter_grid_stat{})
	db.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb")
	db.Exec("SELECT create_hypertable('meter_grid_stats', 'timestamp')")
}

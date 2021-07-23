package server

import (
	"fmt"
	"os"
	"strings"
	"time"

	libs "github.com/lhxlnsy/paplibs"

	"github.com/gomodule/redigo/redis"
)

var Redis = StartRedis()

type RedisPool interface {
	Ready()
	PushData(key string, value interface{}) error
	GetData(key string) []string
	Set(key string, val string) error
	Get(key string) (string, error)
	GetLen(key string) (int, error)
	ping()
	testStore()
}

type PAPRedis struct {
	pool       *redis.Pool
	listlength int
}

func StartRedis() *PAPRedis {
	redispool := initPool()
	return &PAPRedis{
		pool:       redispool,
		listlength: 60,
	}
}

func (r *PAPRedis) GetTimeStamp() time.Time {
	conn := r.pool.Get()
	defer conn.Close()
	s, _ := redis.String(conn.Do("LINDEX", "Timestamp", "0"))
	return libs.StrToTime(s)
}

func (r *PAPRedis) GetDefaultLen() int {
	return r.listlength
}

func (r *PAPRedis) GetLen(key string) (int, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("LLEN", key))
}

func (r *PAPRedis) EmptyList(key string) {
	conn := r.pool.Get()
	defer conn.Close()
	//TODO: EMPTY THE LIST AFTER WE GET THE DATA
	_, err := conn.Do("DEL", key)
	if err != nil {
		libs.Logger.Print().Errorf("ERROR: fail to trim key %s, error %s", key, err.Error())
	}
}

func (r *PAPRedis) GetData(key string) []string {
	conn := r.pool.Get()
	defer conn.Close()

	s, err := redis.Strings(conn.Do("LRANGE", key, "0", "-1"))
	if err != nil {
		libs.Logger.Print().Errorf("ERROR: fail get key %s, error %s", key, err.Error())
	}
	libs.Logger.Print().Infof("return data for key :", key)
	libs.Logger.Print().Info(s)
	return s
}
func initPool() *redis.Pool {
	//libs.PAPLogger.Log.Sugar().Infof(libs.Config.Redis.Host)
	pool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", libs.Config.Redis.Host+":"+libs.Config.Redis.Port)
			if err != nil {
				libs.Logger.Print().Errorf("ERROR: fail init redis: %s \n", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
	return pool
}

func (r *PAPRedis) PushData(key string, value interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	var redisstr string
	switch value.(type) {
	case float64:
		redisstr = fmt.Sprintf("%f", value)
	default:
		redisstr = fmt.Sprintf("%v", value)
	}
	_, err := conn.Do("LPUSH", key, redisstr)
	if err != nil {
		libs.Logger.Print().Errorf("ERROR: fail set key %s, val %v, error %s", key, redisstr, err.Error())
		return err
	}
	_, err = conn.Do("LTRIM", key, 0, r.listlength-1)
	if err != nil {
		libs.Logger.Print().Errorf("ERROR: fail trim the list to length %d, error %s", r.listlength, err.Error())
		return err
	}
	return nil
}

func (r *PAPRedis) ping() {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		libs.Logger.Print().Errorf("ERROR: fail ping redis conn: %s", err.Error())
		os.Exit(1)
	}
}

func (r *PAPRedis) Set(key string, val string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		libs.Logger.Print().Errorf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
		return err
	}
	libs.Logger.Print().Infof("Done: Set Redis Key: %s, value: %s", key, val)
	return nil
}

func (r *PAPRedis) Get(key string) (string, error) {
	// get conn and put back when exit from method
	conn := r.pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		libs.Logger.Print().Errorf("ERROR: fail get key %s, error %s", key, err.Error())
		return "", err
	}
	return s, nil
}

func (r *PAPRedis) testStore() {
	// get conn and put back when exit from method
	conn := r.pool.Get()
	defer conn.Close()

	macs := []string{"Version||1.0", "Author||Jeffrey Li", "CopyRight||Planet Ark Power"}
	for _, mac := range macs {
		pair := strings.Split(mac, "||")
		r.Set(pair[0], pair[1])
	}
}

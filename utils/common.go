package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func RedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   3,
		MaxActive: 20, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func ScheduleTask(queue string, data string, delay int64) {
	c := RedisPool().Get()
	defer c.Close()

	run_at := time.Now().Unix() + delay + 1

	c.Do("ZADD", queue, run_at, data)

	return
}

package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GetMD5Hash returns MD5 hash for a given string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// func RedisPool() *redis.Pool {
// 	return &redis.Pool{
// 		MaxIdle:   3,
// 		MaxActive: 20, // max number of connections
// 		Dial: func() (redis.Conn, error) {
// 			c, err := redis.Dial("tcp", ":6379")
// 			if err != nil {
// 				panic(err.Error())
// 			}
// 			return c, err
// 		},
// 	}
// }

// func ScheduleTask(queue string, data string, delay int64) {
// 	c := RedisPool().Get()
// 	defer c.Close()

// 	run_at := time.Now().Unix() + delay + 1

// 	c.Do("ZADD", queue, run_at, data)

// 	return
// }

// GetUUID returns a UUID
func GetUUID() string {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

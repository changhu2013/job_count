package main

import (
	"fmt"

	"flag"

	"github.com/garyburd/redigo/redis"
)

func getRedisConnect(redisURL string) redis.Conn {
	var conn, err = redis.DialURL(redisURL)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return conn
}

func jobCountByKey(conn redis.Conn, key string) int64 {
	var count, err = redis.Int64(conn.Do("ZCARD", key))
	if err != nil {
		count = 0
	}

	return count
}

//取客户端总数
func getJobCount(redisURL string, vvv bool) int64 {
	var conn = getRedisConnect(redisURL)
	if conn == nil {
		return 0
	}

	var c0 = jobCountByKey(conn, "mpp_0:jobs")
	var c1 = jobCountByKey(conn, "mpp_1:jobs")

	if vvv {
		fmt.Printf("redisurl[%s] count[%d]", redisURL, c0)
		fmt.Println()
	}

	defer conn.Close()
	return c0 + c1
}

func jobCount(vvv bool, redisURLS []string) int64 {
	var count int64

	for idx := range redisURLS {
		var redisURL = redisURLS[idx]
		count = count + getJobCount(redisURL, vvv)
	}

	return count
}

func main() {

	var v bool
	var t bool

	flag.BoolVar(&v, "v", false, "显示详细信息")
	flag.BoolVar(&t, "t", false, "测试环境")
	flag.Parse()

	var redisURLS = []string{"redis://192.168.200.50:6379"}

	if !t {
		redisURLS = []string{
			"redis://192.168.1.152:6379",
			"redis://192.168.1.152:6380",
			"redis://192.168.1.152:6381",
			"redis://192.168.1.153:6379",
			"redis://192.168.1.153:6380",
			"redis://192.168.1.153:6381",
		}
	}

	var count = jobCount(v, redisURLS)

	fmt.Println(count)
}

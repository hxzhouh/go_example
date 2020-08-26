/*
@Time : 2020/4/7 11:44
@Author : ZhouHui2
@File : main
@Software: GoLand
*/
package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

func main() {

	var c int64
	flag.Int64Var(&c, "c", 21299, "corpId")
	flag.Parse()

	fmt.Println("hello,world")

	client := redis.NewClient(&redis.Options{
		Addr:     "r-bp142699b6b62f14267.redis.rds.aliyuncs.com:6379",
		Password: "r-bp142699b6b62f14:BATy8M24xJ", // no password set
		DB:       0,                               // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("连接redis失败", err)
	}
	result := client.Set(fmt.Sprintf("%s:%d", "OCR_CORPID", c), 1, -1)
	if result.Err() != nil {
		log.Fatal("更新失败", err)
	} else {
		log.Println("success")
	}
}

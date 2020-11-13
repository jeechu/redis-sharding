package main

import (
	"context"
	"fmt"
	"github.com/dunzoit/dunzo_commons/go_commons/common-utils/id_generator"
	"redis-library/redis"
)



func test() {
	redisConn := redis.NewClient()

	for i := 0; i<1000; i++ {
		key := id_generator.GetUniqId()
		redisConn.Set(context.Background(), key, "testvalue", 0)
		fmt.Println("key: ", key)

		result := redisConn.Get(context.Background(), key)
		if result.Err() != nil {
			fmt.Println(result.Err().Error())
			return
		}

		value, err := result.Result()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(value)
	}
}


func main() {
	test()
}

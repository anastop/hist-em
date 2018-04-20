package main

import (
	"net/http"
	"os"
	"github.com/go-redis/redis"
	"time"
	"log"
	"fmt"
        "encoding/json"
	"io/ioutil"
)


func writetoRedis(rediskey string , redisVal string){
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println(redisVal)

	llength, _ := client.LLen(rediskey).Result()
    fmt.Println(llength)
	if llength == 100 {
		client.LPop(rediskey)
	}
	client.RPush(rediskey, redisVal)
}


func main() {
	user := os.Args[1]
	host := os.Args[2]
	netService := os.Args[3]
	app := os.Args[4]
	time.Sleep(time.Second*5)
	for {
		collect(user, host, netService, app)
		time.Sleep(time.Second*1)
	}
}


type CollectorData struct {
  Timestamp int64        `json:"timestamp"`
  Data      json.RawMessage  `json:"data"`
}

func collect(user, host, service, app string) {
	url := "http://167.99.213.33:9000/v1/data"
	resp, err := http.Get(url)
	if err != nil {
          log.Fatal(err)
	}
	responsedata, errori := ioutil.ReadAll(resp.Body)
	if errori != nil {
	  log.Fatal(errori)
	}
	rediskey := fmt.Sprintf("heatmaps:%s:%s:%s:%s:latency", user, host, service, app)
    //b, erro := json.Marshal(resp)
	fmt.Println(string(responsedata))

        m := CollectorData{Data: responsedata, Timestamp: time.Now().Unix()}
        data, errorios := json.Marshal(m)
        if errorios != nil {
          log.Fatal(errorios)
        }
        writetoRedis(rediskey, string(data))

}


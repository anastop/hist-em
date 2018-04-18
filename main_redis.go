package main

import (
//  "net/http"
  "strings"
  "encoding/json"
  "github.com/go-redis/redis"
  "time"
  "os/exec"
  "log"
  "fmt"
  "regexp"
  "bufio"
)


type histogramColumn struct {
    Value  string `json:"value"`
    Count  string `json:"count"`
}

type histogram struct {
    Histogram []histogramColumn `json:"histogram"`
}


func writetoRedis(redisVal histogram){
  client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
  })

  val, err := client.Get("agg_data").Result()
  if err != nil {
    panic(err)
  }
  fmt.Println("key", val)
  var dat []histogram
  bytes := []byte(val)
  if err := json.Unmarshal(bytes, &dat); err != nil {
        panic(err)
    }
  fmt.Println(dat)
  fmt.Println("len:", len(dat))
  var redisValArr []histogram
  if len(dat) == 5 {
    redisValArr = append(dat[1:], redisVal)
  }else{
    redisValArr = append(dat, redisVal)
  }
  js, err := json.Marshal(redisValArr)
  if err != nil {
    panic(err)
    return
  }

  err = client.Set("agg_data", js, 0).Err()
  if err != nil {
    panic(err)
  }
}

func main() {
   c := make(chan histogram)
   for {
	execute_com()
        time.Sleep(time.Second*5)
	}
}

func execute_com() {

  // var wg sync.WaitGroup
  // go esDeamon()
  cmd := exec.Command("/home/nfvacc/bcc/tools/runqlat.py", "2", "1")
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    log.Fatal(err)
  }


  in := bufio.NewScanner(stdout)
  err = cmd.Start()
  if err != nil {
    log.Fatal(err)
  }

  count := 0
  var  histArray histogram
  var histArr []histogramColumn
  for in.Scan() {
      line := in.Text()
      count++
      if count > 3 {
          s := strings.Split(line, "|")
          s2 := strings.Split(s[0], ":")

          s3 := strings.Split(s2[0], "->")
          re := regexp.MustCompile("[0-9]+")
          value := re.FindAllString(s3[1], -1)[0]

          count := re.FindAllString(s2[1], -1)[0]
          fmt.Println(value, count)
          histArr = append(histArr, histogramColumn{Value: value, Count: count})
      }
  }
  histArray.Histogram = histArr
  fmt.Println(histArray)
  writetoRedis(histArray)
}

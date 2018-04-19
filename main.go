package main

import (
  "net/http"
  "strings"
  "encoding/json"
  "github.com/cleversoap/go-cp"
  "os/exec"
  "log"
  "fmt"
  "bufio"
)

var quit chan struct{}

type histogramColumn struct {
    Value  string `json:"value"`
    Count  string `json:"count"`
}

type histogram struct {
    Histogram []histogramColumn `json:"histogram"`
}




func run_cyclictest() *exec.Cmd{
  cmd := exec.Command("sh", "-c", "./cyclictest -D3 -h100 > outputGo")
  err := cmd.Start()
  if err != nil {
    log.Fatal(err)
  }
  return cmd
}

func start_cyclictest() {
  var cmd *exec.Cmd
  for {
    select {
    case <-quit:
      // kill_cyclictest(cmd)
      break
    default:
      cmd = run_cyclictest()
      cmd.Wait()
      err := cp.Copy("outputGo", "output.txt")
      if err != nil {
        panic(err)
      }
    }
  }

}

func kill_cyclictest(cmd *exec.Cmd){
  // Kill it:
  if err := cmd.Process.Kill(); err != nil {
    log.Fatal("failed to kill process: ", err)
  }
}


func parse_cyclictest_results() histogram{
  cmd := exec.Command("sh", "-c", "python ./create_hist.py")
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    log.Fatal(err)
  }

  in := bufio.NewScanner(stdout)

  err = cmd.Start()
  if err != nil {
    log.Fatal(err)
  }
  var  histArray histogram
  var histArr []histogramColumn
  for in.Scan(){
    line := in.Text()
    fmt.Println(line)
    s := strings.Split(line, " ")
    histArr = append(histArr, histogramColumn{Value: s[0], Count: s[1]})
  }
  histArray.Histogram = histArr
  return histArray
}

func stopCyclRoutine(){
  quit <- struct{}{}
}



func HTTPDataHandler(w http.ResponseWriter, r *http.Request) {
  ret := parse_cyclictest_results()
  js, err := json.Marshal(ret)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)

}

func HTTPStopHandler(w http.ResponseWriter, r *http.Request) {
   stopCyclRoutine()
   w.Header().Set("Server", "Cyclic Web Server")
   w.WriteHeader(200)
}


func HTTPStartHandler(w http.ResponseWriter, r *http.Request) {

   quit = make(chan struct{})
   go start_cyclictest()
   w.Header().Set("Server", "Cyclic Web Server")
   w.WriteHeader(200)

}


func main() {

   http.HandleFunc("/v1/data", HTTPDataHandler)
   http.HandleFunc("/v1/start", HTTPStartHandler)
   http.HandleFunc("/v1/stop", HTTPStopHandler)
   log.Fatal(http.ListenAndServe(":9000", nil))
}




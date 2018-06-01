package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cleversoap/go-cp"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

var quit chan struct{}

type histogramColumn struct {
	Value float64 `json:"value"`
	Count int64   `json:"count"`
}

type histogram struct {
	Histogram []histogramColumn `json:"histogram"`
}

func runCyclictest() *exec.Cmd {
	cmd := exec.Command("sh", "-c", "./cyclictest -q -D2 -h100 > outputGo")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	return cmd
}

func startCyclictest() {
	var cmd *exec.Cmd
	for {
		select {
		case <-quit:
			// kill_cyclictest(cmd)
			break
		default:
			cmd = runCyclictest()
			cmd.Wait()
			err := cp.Copy("outputGo", "output.txt")
			if err != nil {
				panic(err)
			}
		}
	}

}

func killCyclictest(cmd *exec.Cmd) {
	// Kill it:
	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
	}
}

func parseCyclictestResults() histogram {
	var val float64
	var count int64
	var errori error
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
	var histArray histogram
	var histArr []histogramColumn
	for in.Scan() {
		line := in.Text()
		fmt.Println(line)
		s := strings.Split(line, " ")
		if val, errori = strconv.ParseFloat(s[0], 64); errori != nil {
			log.Fatal(errori)
		}
		if count, errori = strconv.ParseInt(s[1], 10, 64); errori != nil {
			log.Fatal(errori)
		}
		col := histogramColumn{
			Value: val,
			Count: count,
		}
		histArr = append(histArr, col)
	}
	histArray.Histogram = histArr
	fmt.Println()
	return histArray
}

func stopCyclRoutine() {
	quit <- struct{}{}
}

func HTTPDataHandler(w http.ResponseWriter, r *http.Request) {
	ret := parseCyclictestResults()
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
	go startCyclictest()
	w.Header().Set("Server", "Cyclic Web Server")
	w.WriteHeader(200)

}

func main() {
	go startCyclictest()
	http.HandleFunc("/v1/data", HTTPDataHandler)
	http.HandleFunc("/v1/start", HTTPStartHandler)
	http.HandleFunc("/v1/stop", HTTPStopHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

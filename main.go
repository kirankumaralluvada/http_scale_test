package main

import (
	"fmt"
	"time"
	"net/http"
	"os"
	_ "strconv"
	_ "math/rand"

	//vegeta "github.com/tsenart/vegeta/lib"
	vegeta "github.com/kirankumaralluvada/vegeta/lib"
)

func main() {
	fmt.Printf("Starting tests\n")
	rate := vegeta.Rate{Freq: 20, Per: time.Second}
	duration := 2 * time.Second
	base_target := vegeta.Target{
				Method: "GET",
				URL:    "http://10.106.40.189:8080/foo",
				Header: http.Header{
				//	"x-nws-device-id":     []string{"135015946"},
					"x-nws-sequence-no":   []string{"1"},
					"Host":                []string{"10.46.2.5"},
				},
			}

/*
	list_target := []vegeta.Target{}

	for id := 2; id<10; id++ {
		base_target.Header["x-nws-device-id"] = []string{strconv.Itoa(id)}
		list_target = append(list_target, base_target)
	}
*/

	targeter := vegeta.NewStaticTargeter(base_target)
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		fmt.Println(res.Latency)
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Sent: %d, Success: %f\n", metrics.Requests, metrics.Success)
	fmt.Printf("Max: %s, Min: %s, Mean: %s, 99th percentile: %s\n", metrics.Latencies.Max, metrics.Latencies.Min, metrics.Latencies.Mean, metrics.Latencies.P99)

	rep := vegeta.NewHDRHistogramPlotReporter(&metrics)

	filename, _ := os.Create("out.txt")
	rep.Report(filename)
	filename.Close()
}

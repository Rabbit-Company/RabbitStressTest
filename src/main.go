package main

import (
	"flag"
	"fmt"
	"github.com/aquasecurity/table"
	"github.com/valyala/fasthttp"
	"os"
	"sync"
	"time"
)

// Colors
var reset string = "\033[0m"
var green string = "\033[32m"

const maxInt int = int(^uint(0) >> 1)

var target string
var req int
var duration int
var delay int
var workers int

var success int = 0
var errors int = 0

var deliveryTimes []int64 = make([]int64, 0)

func init() {
	flag.StringVar(&target, "u", "http://127.0.0.1", "Specify URL / target.")
	flag.IntVar(&req, "r", 100, "Specify number of requests.")
	flag.IntVar(&duration, "t", 0, "Specify duration in seconds.")
	flag.IntVar(&delay, "d", 1, "Specify delay in milliseconds.")
	flag.IntVar(&workers, "w", 10, "Specify number of workers per routine.")
	flag.Parse()
}

func main() {
	fmt.Println(green)
	fmt.Println("██████╗  █████╗ ██████╗ ██████╗ ██╗████████╗    ███████╗████████╗██████╗ ███████╗███████╗███████╗    ████████╗███████╗███████╗████████╗")
	fmt.Println("██╔══██╗██╔══██╗██╔══██╗██╔══██╗██║╚══██╔══╝    ██╔════╝╚══██╔══╝██╔══██╗██╔════╝██╔════╝██╔════╝    ╚══██╔══╝██╔════╝██╔════╝╚══██╔══╝")
	fmt.Println("██████╔╝███████║██████╔╝██████╔╝██║   ██║       ███████╗   ██║   ██████╔╝█████╗  ███████╗███████╗       ██║   █████╗  ███████╗   ██║   ")
	fmt.Println("██╔══██╗██╔══██║██╔══██╗██╔══██╗██║   ██║       ╚════██║   ██║   ██╔══██╗██╔══╝  ╚════██║╚════██║       ██║   ██╔══╝  ╚════██║   ██║   ")
	fmt.Println("██║  ██║██║  ██║██████╔╝██████╔╝██║   ██║       ███████║   ██║   ██║  ██║███████╗███████║███████║       ██║   ███████╗███████║   ██║   ")
	fmt.Println("╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚═════╝ ╚═╝   ╚═╝       ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝       ╚═╝   ╚══════╝╚══════╝   ╚═╝   ")
	fmt.Println(reset)

	//Validation
	if req < 1 {
		req = 100
	}
	if duration < 0 {
		duration = 0
	}
	if delay < 0 {
		delay = 0
	}
	if workers > req {
		workers = req
	}

	t := table.New(os.Stdout)
	t.SetLineStyle(table.StyleBlue)
	t.SetAlignment(table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignCenter)

	if duration != 0 {
		req = maxInt
		t.SetHeaders("Target", "Duration", "Workers", "Delay")
		t.AddRow(target, fmt.Sprintf("%ds", duration), fmt.Sprintf("%d", workers), fmt.Sprintf("%dms", delay))
	} else {
		t.SetHeaders("Target", "Requests", "Workers", "Delay")
		t.AddRow(target, fmt.Sprintf("%d", req), fmt.Sprintf("%d", workers), fmt.Sprintf("%dms", delay))
	}
	t.Render()
	fmt.Println("")

	var wg sync.WaitGroup

	var client = &fasthttp.Client{
		MaxConnsPerHost:     maxInt,
		MaxIdleConnDuration: 500,
		Dial: (&fasthttp.TCPDialer{
			Concurrency: maxInt,
		}).Dial,
	}

	start := time.Now()
	for i := 1; i <= req/workers; i++ {
		if duration != 0 {
			if time.Since(start).Seconds() >= float64(duration) {
				break
			}
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
		wg.Add(1)
		go func() {
			defer wg.Done()
			var body []byte
			for j := 1; j <= workers; j++ {
				start2 := time.Now()
				status, _, _ := client.Get(body, target)
				if status == 200 {
					success++
					deliveryTimes = append(deliveryTimes, time.Since(start2).Milliseconds())
				} else {
					errors++
				}
			}
		}()
	}
	wg.Wait()
	secs := time.Since(start).Seconds()

	errorRate := (float32(errors) / float32(success+errors) * 100)
	var total int64 = 0
	for _, number := range deliveryTimes {
		total = total + number
	}
	average := total / int64(len(deliveryTimes))

	t = table.New(os.Stdout)
	if errors == 0 {
		t.SetLineStyle(table.StyleGreen)
	} else if errorRate < 50 {
		t.SetLineStyle(table.StyleYellow)
	} else {
		t.SetLineStyle(table.StyleRed)
	}

	t.SetHeaders("Requests", "Time", "Success", "Errors", "Error rate", "Requests per second", "Time per request")
	t.SetAlignment(table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignCenter)

	t.AddRow(fmt.Sprintf("%d", success+errors), fmt.Sprintf("%.2fs", secs), fmt.Sprintf("%d", success), fmt.Sprintf("%d", errors), fmt.Sprintf("%.2f%%", errorRate), fmt.Sprintf("%.2f", float32(success+errors)/float32(secs)), fmt.Sprintf("%dms", average))
	t.Render()
	fmt.Println("")
}

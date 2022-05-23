package main

import (
	"flag"
	"fmt"
	"github.com/aquasecurity/table"
	"github.com/valyala/fasthttp"
	"github.com/wcharczuk/go-chart/v2"
	"os"
	"sync"
	"time"
)

// Colors
var reset string = "\033[0m"
var green string = "\033[32m"
var blue string = "\033[34m"

const maxInt int = int(^uint(0) >> 1)

var target string
var req int
var duration int
var delay int
var workers int
var graph bool

var success int = 0
var errors int = 0

var deliveryTimes []int64 = make([]int64, 0)

func init() {
	flag.StringVar(&target, "u", "http://127.0.0.1", "Specify URL / target.")
	flag.IntVar(&req, "r", 100, "Specify number of requests.")
	flag.IntVar(&duration, "t", 0, "Specify duration in seconds.")
	flag.IntVar(&delay, "d", 1, "Specify delay in milliseconds.")
	flag.IntVar(&workers, "w", 10, "Specify number of workers per routine.")
	flag.BoolVar(&graph, "g", false, "Create a graph with response times and store it as stats.png")
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

	fmt.Println(blue)
	fmt.Println("┬┌┐┌┌─┐┬ ┬┌┬┐")
	fmt.Println("││││├─┘│ │ │ ")
	fmt.Println("┴┘└┘┴  └─┘ ┴ ")
	fmt.Printf(reset)

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
	fmt.Println(green)
	fmt.Println("┌─┐┬ ┬┌┬┐┌─┐┬ ┬┌┬┐")
	fmt.Println("│ ││ │ │ ├─┘│ │ │ ")
	fmt.Println("└─┘└─┘ ┴ ┴  └─┘ ┴ ")
	fmt.Printf(reset)

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
	var slowest int64 = 0
	var fastest int64 = 86400000
	for _, number := range deliveryTimes {
		if number > slowest {
			slowest = number
		}
		if number < fastest {
			fastest = number
		}
		total = total + number
	}
	average := total
	if len(deliveryTimes) != 0 {
		average = total / int64(len(deliveryTimes))
	}

	t = table.New(os.Stdout)
	if errors == 0 {
		t.SetLineStyle(table.StyleGreen)
	} else if errorRate < 50 {
		t.SetLineStyle(table.StyleYellow)
	} else {
		t.SetLineStyle(table.StyleRed)
	}

	t.SetHeaders("Requests", "Time", "Success", "Errors", "Error rate", "Requests per second")
	t.SetAlignment(table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignCenter, table.AlignCenter)

	t.AddRow(fmt.Sprintf("%d", success+errors), fmt.Sprintf("%.2fs", secs), fmt.Sprintf("%d", success), fmt.Sprintf("%d", errors), fmt.Sprintf("%.2f%%", errorRate), fmt.Sprintf("%.2f", float32(success+errors)/float32(secs)))
	t.Render()

	fmt.Println(green)
	fmt.Println("┬─┐┌─┐┌─┐ ┬ ┬┌─┐┌─┐┌┬┐┌─┐")
	fmt.Println("├┬┘├┤ │─┼┐│ │├┤ └─┐ │ └─┐")
	fmt.Println("┴└─└─┘└─┘└└─┘└─┘└─┘ ┴ └─┘")
	fmt.Printf(reset)

	t = table.New(os.Stdout)
	if errors == 0 {
		t.SetLineStyle(table.StyleGreen)
	} else if errorRate < 50 {
		t.SetLineStyle(table.StyleYellow)
	} else {
		t.SetLineStyle(table.StyleRed)
	}

	t.SetHeaders("Fastest", "Average", "Slowest")
	t.SetAlignment(table.AlignCenter, table.AlignCenter, table.AlignCenter)

	t.AddRow(fmt.Sprintf("%dms", fastest), fmt.Sprintf("%dms", average), fmt.Sprintf("%dms", slowest))
	t.Render()
	fmt.Println("")

	if graph {
		createGraph()
	}
}

func createGraph() {
	stats := make([]float64, len(deliveryTimes))
	stats2 := make([]float64, len(deliveryTimes))

	for i := range deliveryTimes {
		stats2[i] = float64(i)
		stats[i] = float64(deliveryTimes[i])
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Request",
			Style: chart.Style{
				TextRotationDegrees: 45,
			},
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d", int(v.(float64)))
			},
		},
		YAxis: chart.YAxis{
			Name: "Latency",
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d ms", int(v.(float64)))
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				Name:    target,
				XValues: stats2,
				YValues: stats,
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, _ := os.Create("stats.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}

package main

import (
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"sync"
	"time"
)

// Colors
var reset string = "\033[0m"
var red string = "\033[31m"
var green string = "\033[32m"
var blue string = "\033[34m"

const maxInt int = int(^uint(0) >> 1)

var target string
var req int
var duration int
var delay int

var success int = 0
var errors int = 0

func init() {
	flag.StringVar(&target, "t", "http://127.0.0.1", "Specify target / url.")
	flag.IntVar(&req, "r", 100, "Specify number of requests.")
	flag.IntVar(&duration, "d", 0, "Specify duration in seconds.")
	flag.IntVar(&delay, "w", 1, "Specify delay in milliseconds.")
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

	fmt.Println("Target: " + blue + target + reset)
	if duration != 0 {
		req = maxInt
		fmt.Printf("Duration: %s%ds%s\n", blue, duration, reset)
	} else {
		fmt.Printf("Requests: %s%d%s\n", blue, req, reset)
	}
	fmt.Printf("Delay: %s%dms%s\n", blue, delay, reset)

	var wg sync.WaitGroup

	var client = &fasthttp.Client{
		MaxConnsPerHost:     maxInt,
		MaxIdleConnDuration: 500,
		Dial: (&fasthttp.TCPDialer{
			Concurrency: maxInt,
		}).Dial,
	}

	start := time.Now()
	for i := 1; i <= req; i++ {
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
			status, _, _ := client.Get(body, target)
			if status == 200 {
				success++
			} else {
				errors++
			}
		}()
	}
	wg.Wait()
	secs := time.Since(start).Seconds()
	fmt.Printf("\n%s%d%s requests where performed in %s%.2fs%s\n", green, success+errors, reset, green, secs, reset)
	fmt.Printf("---------------------\n")
	fmt.Printf("Success: %s%d%s\n", green, success, reset)
	if errors == 0 {
		fmt.Printf("Errors: %s%d%s\n", green, errors, reset)
		fmt.Printf("---------------------\n")
		fmt.Printf("Error rate: %s%.2f%%%s\n", green, (float32(errors) / float32(success+errors) * 100), reset)
	} else {
		fmt.Printf("Errors: %s%d%s\n", red, errors, reset)
		fmt.Printf("---------------------\n")
		fmt.Printf("Error rate: %s%.2f%%%s\n", red, (float32(errors) / float32(success+errors) * 100), reset)
	}
	fmt.Printf("---------------------\n")
	fmt.Printf("Requests per second: %s%.2f%s\n", green, float32(success+errors)/float32(secs), reset)

	fmt.Println(reset)
}

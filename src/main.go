package main

import (
	"fmt"
	"flag"
	"time"
	"sync"
	"github.com/valyala/fasthttp"
)

// Colors
var reset string = "\033[0m"
var red string = "\033[31m"
var green string = "\033[32m"
var blue string = "\033[34m"

var target string
var req int
var duration int

var success int = 0
var errors int = 0

func init() {
	flag.StringVar(&target, "t", "https://google.com", "Specify target / url.")
	flag.IntVar(&req, "r", 100, "Specify number of requests.")
	flag.IntVar(&duration, "d", 10, "Specify duration in seconds.")
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
	fmt.Printf("Requests: %s%d%s\n", blue, req, reset)
	fmt.Printf("Duration: %s%ds%s\n\n", blue, duration, reset)

	var wg sync.WaitGroup
	
	var client = &fasthttp.Client{
		MaxConnsPerHost: req*2,
	}
	
	start := time.Now()
	for i:=1; i <= req; i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			var body []byte
			status, _, _ := client.Get(body, target)
			if(status == 200){
				success++;
			}else{
				errors++;
			}
		}()
	}
	wg.Wait()
	secs := time.Since(start).Seconds()
	fmt.Printf("%s%d%s requests where performed in %s%fs%s\n", green, success+errors, reset, green, secs, reset)
	fmt.Printf("Success: %s%d%s\n", green, success, reset)
	if(errors == 0){
		fmt.Printf("Errors: %s%d%s\n", green, errors, reset)
	}else{
		fmt.Printf("Errors: %s%d%s\n", red, errors, reset)
	}

	fmt.Println(reset)
}
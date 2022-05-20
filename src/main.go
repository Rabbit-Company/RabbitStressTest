package main

import (
	"os"
	"fmt"
	"flag"
)

// Colors
var reset string = "\033[0m"
var red string = "\033[31m"
var green string = "\033[32m"
var blue string = "\033[34m"

var target string = "google.com"
var reqps int = 100
var duration int = 10

func main() {
	fmt.Println(green)
	fmt.Println("██████╗  █████╗ ██████╗ ██████╗ ██╗████████╗    ███████╗████████╗██████╗ ███████╗███████╗███████╗    ████████╗███████╗███████╗████████╗")
	fmt.Println("██╔══██╗██╔══██╗██╔══██╗██╔══██╗██║╚══██╔══╝    ██╔════╝╚══██╔══╝██╔══██╗██╔════╝██╔════╝██╔════╝    ╚══██╔══╝██╔════╝██╔════╝╚══██╔══╝")
	fmt.Println("██████╔╝███████║██████╔╝██████╔╝██║   ██║       ███████╗   ██║   ██████╔╝█████╗  ███████╗███████╗       ██║   █████╗  ███████╗   ██║   ")
	fmt.Println("██╔══██╗██╔══██║██╔══██╗██╔══██╗██║   ██║       ╚════██║   ██║   ██╔══██╗██╔══╝  ╚════██║╚════██║       ██║   ██╔══╝  ╚════██║   ██║   ")
	fmt.Println("██║  ██║██║  ██║██████╔╝██████╔╝██║   ██║       ███████║   ██║   ██║  ██║███████╗███████║███████║       ██║   ███████╗███████║   ██║   ")
	fmt.Println("╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚═════╝ ╚═╝   ╚═╝       ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝       ╚═╝   ╚══════╝╚══════╝   ╚═╝   ")
	fmt.Println(reset)

	if(len(os.Args) != 1){
		target = os.Args[1]
	}

	flag.IntVar(&reqps, "rps", 100, "Specify Requests per second.")
	flag.IntVar(&duration, "duration", 10, "Specify duration in seconds.")
	flag.Parse()

	fmt.Println("Target: " + blue + target + reset)	
	fmt.Printf("Requests per sec: %s%d%s\n", blue, reqps, reset)
	fmt.Printf("Duration: %s%ds%s\n", blue, duration, reset)

	fmt.Println(reset)
}
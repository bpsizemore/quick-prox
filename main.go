package main

import (
	"net/http"
	"log"
	"github.com/elazarl/goproxy"
	"flag"
	"fmt"
	"io/ioutil"
)

var log_type string

func main() {

	/* Goal: A simple utility for quickly and easily logging network requests from 
		command line tools. In particular, this should serve as a simpler alternative
		to running burp or other similar proxies just to get requests logs
		from a CLI application.

	   Use cases:
		* Debugging scripts that make network requests
		* Investigating what network requests are made by binaries/scripts

	   TODOs:
	   	* Add response information to the log (e.g. server response status code.)
		* Add option to log to file
		
	*/

	//stdoutFlag := flag.Bool("stdout", false, "Print logs to stdout.")
	detailedFlag := flag.Bool("detailed", false, "Apply detailed logging to all requests.")
	veryDetailedFlag := flag.Bool("very", false, "Apply very detailed logging to all requests.")
	portFlag := flag.String("interface", ":8080", "proxy port, default is :8080 indicating all interfaces on port 8080")

	flag.Parse()

	if *veryDetailedFlag == true {
		log_type = "very detailed"
	} else if *detailedFlag == true {
		log_type = "detailed"
	} else {
		log_type = "basic"
	}
	/*if *stdoutFlag == true {
		fmt.Println("Standard out mode selected, printing logs to terminal...")
	}*/

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx)(*http.Request, *http.Response) {
			switch log_type {
			case "very detailed":
				log.Println(r.Method, r.URL)
				for header := range r.Header {
					fmt.Println("	",header,r.Header[header])
				}
				body, _ := ioutil.ReadAll(r.Body)
				fmt.Println(string(body))
			case "detailed":
				log.Println(r.Method, r.URL)
				for header := range r.Header {
					fmt.Println("	",header,r.Header[header])
				}
			default:
				log.Println(r.Method, "request sent to", r.Host)
			}
			fmt.Println()
			return r,nil
		})

	log.Fatal(http.ListenAndServe(*portFlag, proxy))
}

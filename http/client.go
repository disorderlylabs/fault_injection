package main

import (
	"os"
	"fmt"
	"flag"
	"strconv"
	"net/http"
)

func inject_fault(serviceName string, faultType string) {



	//if a serviceName is entered, then we must be targeting it for a fault injection
	if serviceName != "" {
		switch faultType {
		case "delay_ms":
			time, _ := strconv.ParseInt(flag.Args()[0], 10, 64)
			fmt.Printf("Fault type: delay of %d ms to service %s\n", time, serviceName)
		case "http_error":
			errCode, _ := strconv.ParseInt(flag.Args()[0], 10, 64)
			fmt.Printf("Injecting http error code %d into service %s\n", errCode, serviceName)

		case "drop_packet":
			fmt.Printf("Fault type: dropping packet going to service %s\n", serviceName)
		default:
			fmt.Fprintf(os.Stderr, "error: must specify fault type for the service\n")
			os.Exit(1)

		}
	}

	req, _ := http.NewRequest("GET", "http://localhost:8080/svc1", nil)

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("Error serving client trace request")
	}
}

func dump_trace() {
	//Requesting collector server to dump spans
	req, _ := http.NewRequest("GET", "http://localhost:10000/dump", nil)
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("Error requesting span dump")
	}

}

func main() {
	
	var (
		request = flag.String("request", "", "fault_inject dump_trace")
		serviceName = flag.String("serviceName", "", "service targeted for fault injection")
		faultType = flag.String("faultType", "", "delay_ms, drop_packet")		
	)

	//parse the commandline arguments
	flag.Parse()

	if *request == "fault_inject" {
		inject_fault(*serviceName, *faultType)
	}else if *request == "dump_trace" {
		dump_trace()
	}
	
}
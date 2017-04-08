package main

import (
	"os"
	"fmt"
	"flag"
	"strconv"
	"net/http"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/basictracer-go"
	"fault_injection/trace_recorder/dapperish"
	"io/ioutil"
)

var(
	sp opentracing.Span
	tracer opentracing.Tracer
)


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
		faultType = flag.String("faultType", "", "delay_ms, http_error")
	)

	//parse the commandline arguments
	flag.Parse()


	if *request == "fault_inject" {
		//inject_fault(*serviceName, *faultType)
		tracer = basictracer.New(dapperish.NewTrivialRecorder("fi"))
		opentracing.InitGlobalTracer(tracer)


		sp = opentracing.StartSpan("test_client")
		defer sp.Finish()

		if *serviceName != "" {
			switch *faultType {
				case "delay_ms":
					time, _ := strconv.ParseInt(flag.Args()[0], 10, 64)
					fmt.Printf("Fault type: delay of %d ms to service %s\n", time, *serviceName)
					sp.SetBaggageItem("injectfault", (*serviceName + "_delay:" + flag.Args()[0]))
				case "http_error":
					errCode, _ := strconv.ParseInt(flag.Args()[0], 10, 64)
					fmt.Printf("Injecting http error code %d into service %s\n", errCode, *serviceName)
					sp.SetBaggageItem("injectfault", (*serviceName + "_errcode:" + flag.Args()[0]))
				default:
					fmt.Fprintf(os.Stderr, "error: must specify fault type for the service\n")
				os.Exit(1)

			}
		}

		fmt.Println("request: " + *request + " serviceName: " + *serviceName + " faultType: " + *faultType)

		req, _ := http.NewRequest("GET", "http://localhost:8080/svc1", nil)

		err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			fmt.Printf("failed to inject")
		}
		fmt.Println(sp.BaggageItem("injectfault"))

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Errorf("Request error")
		}
		resp_body, _ := ioutil.ReadAll(response.Body)

		fmt.Println(response.Status)
		fmt.Println(string(resp_body))

	}else if *request == "dump_trace" {
		dump_trace()
	}
	
}
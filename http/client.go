package main

import (
	"os"
	"fmt"
	"flag"
	"strconv"
	"net/http"
	"io/ioutil"
	"fault_injection/trace_recorder/dapperish"
	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
)

func main() {
	
	var (
		sp opentracing.Span
		tracer opentracing.Tracer
		//injectFault bool = false
		serviceName = flag.String("serviceName", "", "service targeted for fault injection")
		faultType = flag.String("faultType", "", "delay_ms, drop_packet")		
	)
	
	//initialize tracer
	tracer = basictracer.New(dapperish.NewTrivialRecorder("fi"))
	opentracing.InitGlobalTracer(tracer)
	
	//start a new span
	sp = opentracing.StartSpan("test_client")
	defer sp.Finish()
	
	//parse the commandline arguments
	flag.Parse()
	
	//if a serviceName is entered, then we must be targeting it for a fault injection
	if *serviceName != "" {
		switch *faultType {
			case "delay_ms":
				time, _ := strconv.ParseInt(flag.Args()[0], 10, 64)
				fmt.Printf("Fault type: delay of %d ms to service %s\n", time, *serviceName)	
			    sp.SetBaggageItem("InjectFault", (*serviceName + "_delay:" + flag.Args()[0]))
			
			case "http_error":
				errCode, _ := strconv.ParseInt(flag.Args()[0], 10, 64)
				fmt.Printf("Injecting http error code %d into service %s\n", errCode, *serviceName)
				sp.SetBaggageItem("InjectFault", (*serviceName + "_errcode:" + flag.Args()[0]))
			
			case "drop_packet":
				fmt.Printf("Fault type: dropping packet going to service %s\n", *serviceName)	
				sp.SetBaggageItem("InjectFault", (*serviceName + "_drop"))
			
			default:
				fmt.Fprintf(os.Stderr, "error: must specify fault type for the service\n")
				os.Exit(1)
			
			//*Note: baggage keys are converted into lower case, so look up "injectfault" not "InjectFault"
		}		
	}	
	
	req, _ := http.NewRequest("GET", "http://localhost:8080/svc1", nil)
	
	//q := req.URL.Query()
	//q.Add("api_key", "key_from_environment_or_flag")
	
	err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		fmt.Printf("fail")
	}
	
	response, err := http.DefaultClient.Do(req)
	
	//defer response.Body.Close()
	resp_body, _ := ioutil.ReadAll(response.Body)

    fmt.Println(response.Status)
    fmt.Println(string(resp_body))
	
	
}
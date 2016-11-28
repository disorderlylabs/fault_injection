package main

import (
	"fmt"
	"flag"
	"strconv"
	"net/http"
	"io/ioutil"
	"fault_injection/dapperish"
	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
)

func main() {
		
	var (
		tracer opentracing.Tracer
		serviceName = flag.String("service", "service1", "service targeted for fault injection")
		faultType = flag.String("inject", "", "delay_ms, dropped_package")		
	)
	
	
	
	switch *faultType {
		case "delay":
		time, _ := strconv.ParseInt(flag.Args()[0], 10, 64)
		fmt.Printf("Inject delay of %d to service %s\n", time, serviceName)
		
	}

	tracer = basictracer.New(dapperish.NewTrivialRecorder("fi"))
	opentracing.InitGlobalTracer(tracer)

	var sp opentracing.Span
	sp = opentracing.StartSpan("test_client")
	defer sp.Finish()
	
	req, _ := http.NewRequest("GET", "http://localhost:8080/svc1", nil)
	
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
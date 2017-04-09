package main

import (
	"net/http"
	"fmt"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	//"temp_fi/core"
	"strconv"
	"os"
	"io/ioutil"
)

const (
	hostPort          = "127.0.0.1:10000"
	collectorEndpoint = "http://localhost:10000/collect"
	sameSpan          = true
	traceID128Bit     = true
)

func testDump() {
	req, err := http.NewRequest("GET", "http://localhost:10000/dump", nil)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("Error requesting span dump")
	}
}


func testInject(serviceName string, faultType string, faultVal string, sp opentracing.Span) opentracing.Span{
	switch faultType {
	case "delay_ms":
		time, _ := strconv.ParseInt(faultVal, 10, 64)
		fmt.Printf("Fault type: delay of %d ms to service %s\n", time, serviceName)
		sp.SetBaggageItem("injectfault", (serviceName + "_delay:" + faultVal))

	case "http_error":
		errcode, _ := strconv.ParseInt(faultVal, 10, 64)
		fmt.Printf("Injecting http error code %d into service %s\n", errcode, serviceName)
		sp.SetBaggageItem("injectfault", (serviceName + "_errcode:" + faultVal))

	default:
		fmt.Fprintf(os.Stderr, "error: must specify fault type for the service\n")
		os.Exit(1)
	}

	return sp
}




func main() {
	collector, err := zipkin.NewHTTPCollector(collectorEndpoint)
	if err != nil {
		fmt.Printf("unable to create a collector: %+v", err)
		os.Exit(-1)
	}

	recorder := zipkin.NewRecorder(collector, false, hostPort, "test")
	// Create our tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
		zipkin.TraceID128Bit(traceID128Bit),
	)

	if err != nil {
		fmt.Printf("unable to create tracer: %+v", err)
		os.Exit(-1)
	}
	opentracing.InitGlobalTracer(tracer)


	sp := opentracing.StartSpan("Test")
	req, _ := http.NewRequest("GET", "http://localhost:8080/service1", nil)



	//TEST: inject 100ms delay on service4
	//sp = testInject("service4", "delay_ms", "100", sp)

	//TEST: inject internal server error
	//sp = testInject("service4", "http_error", "500", sp)



	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		fmt.Printf("failed to inject")
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("Request error")
	}
	resp_body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(response.Status)
	fmt.Println(string(resp_body))



	//*Uncomment to dump traces to JSON file
	testDump()


}





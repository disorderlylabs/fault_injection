package main

import (
	"os"
	"fmt"
	//"flag"
	"strconv"
	"net/http"
	"github.com/opentracing/opentracing-go"
	//"github.com/opentracing/basictracer-go"
	"io/ioutil"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	//"strings"
)

var(
	sp opentracing.Span
	tracer opentracing.Tracer

)

const (
	hostPort          = "127.0.0.1:10000"
	collectorEndpoint = "http://localhost:10000/collect"
	sameSpan          = true
	traceID128Bit     = true
)


func dump_trace() {
	req, err := http.NewRequest("GET", "http://localhost:10000/dump", nil)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("Error requesting span dump")
	}

}

func main() {
	collector, err := zipkin.NewHTTPCollector(collectorEndpoint)
	if err != nil {
		fmt.Printf("unable to create a collector: %+v", err)
		os.Exit(-1)
	}

	recorder := zipkin.NewRecorder(collector, false, hostPort, "client")
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


	sp = opentracing.StartSpan("test_client")
	defer sp.Finish()

	request := "fault_inject"
	//request := "dump_trace"
	serviceName := "service4"
	faultType := "delay_ms"
	faultVal := "10"


	if request == "fault_inject" {
		if serviceName != "" {
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


		}

		req, _ := http.NewRequest("GET", "http://localhost:8080/service1", nil)

		err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(req.Header))
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


	}else if request == "dump_trace" {
		dump_trace()
	}


	
}
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"

	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"fault_injection/trace_recorder/dapperish"
	"fault_injection/http/common"
)


func service1(w http.ResponseWriter, r *http.Request) {
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))
	span := common.Check_and_start_span(err, "SERVICE_1", spCtx)	
	defer span.Finish()
	
	span.SetBaggageItem("svc1_msg", "hello_from_svc1")
	span.SetBaggageItem("svc1_svcs_invoked", "2|3|4")
	
	///Requesting service 2
	svc2_req, _ := http.NewRequest("GET", "http://localhost:8080/svc2", nil)
	
	// Inject the trace information into the HTTP Headers.
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc2_req.Header))
	common.Check_inject_error(err, r)
	
	_, err = http.DefaultClient.Do(svc2_req)
	common.Check_request_error(err, "service_2", r)
		
	
	
	///equesting service 3
	svc3_req, _ := http.NewRequest("GET", "http://localhost:8080/svc3", nil)
	
	// Inject the trace information into the HTTP Headers.
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc3_req.Header))
	common.Check_inject_error(err, r)
	
	_, err = http.DefaultClient.Do(svc3_req)
	common.Check_request_error(err, "service_3", r)
	
	
	
	///Requesting service 4
	svc4_req, _ := http.NewRequest("GET", "http://localhost:8081/svc4", nil)
	
	// Inject the trace information into the HTTP Headers.
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc4_req.Header))
	common.Check_inject_error(err, r)
	
	var response *http.Response
	response, err = http.DefaultClient.Do(svc4_req)
	common.Check_request_error(err, "service_4", r)
	
	resp_body, _ := ioutil.ReadAll(response.Body)

	//fmt.Println("\nResponse status: " + response.Status)
	//fmt.Println("Response body: " + string(resp_body))
	
	fmt.Println()
}//end service1


func service2(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))	
	
	sp = common.Check_and_start_span(err, "SERVICE_2", spCtx)
	defer sp.Finish()
	
	
	sp.LogKV("hello_from", "service_2")
	sp.LogFields(otlog.String("service_2_status", "ok"))
	fmt.Println()
}



func service3(w http.ResponseWriter, r *http.Request) {	
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))
	
	sp = common.Check_and_start_span(err, "SERVICE_3", spCtx)	
	defer sp.Finish()
	
	sp.LogKV("hello_from", "service_3")
	sp.LogFields(otlog.String("service_3_status", "ok"))
	
	sp.SetBaggageItem("svc3_msg", "hello_from_svc3")
	sp.SetBaggageItem("svc3_svcs_invoked", "4")
		
	fmt.Println()
}



func main() {
	
	var tracer opentracing.Tracer
	var port = flag.Int("port", 8080, "Example app port.")

	tracer = basictracer.New(dapperish.NewTrivialRecorder("server_1"))
	opentracing.InitGlobalTracer(tracer)

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/svc1", common.Handler_decorator(service1))
	mux.HandleFunc("/svc2", common.Handler_decorator(service2))
	mux.HandleFunc("/svc3", common.Handler_decorator(service3))
	
	fmt.Printf("Listening on port: %d\n", *port)	
	log.Fatal(http.ListenAndServe(addr, mux))	
}

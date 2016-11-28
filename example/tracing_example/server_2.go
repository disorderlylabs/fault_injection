package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"fault_injection/dapperish"
	"fault_injection/example/tracing_example/common"
)


func service4(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))
	
	
	sp = common.Check_and_start_span(err, "SERVICE_4", spCtx)
	defer sp.Finish()
	
	sp.SetBaggageItem("svc4_msg", "hello_from_svc4")
	sp.SetBaggageItem("svc4_svcs_invoked", "5|6|7")
	
	
	///Requesting service 5
	svc5_req, _ := http.NewRequest("GET", "http://localhost:8081/svc5", nil)
	
	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc5_req.Header))
	common.Check_inject_error(err, r)
	
	_, err = http.DefaultClient.Do(svc5_req)
	common.Check_request_error(err, "service_5", r)
	
	
	
	///Requesting service 6
	svc6_req, _ := http.NewRequest("GET", "http://localhost:8081/svc6", nil)
	
	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc6_req.Header))
	common.Check_inject_error(err, r)
	
	_, err = http.DefaultClient.Do(svc6_req)
	common.Check_request_error(err, "service_6", r)
	
	
	///Requesting service 7
	svc7_req, _ := http.NewRequest("GET", "http://localhost:8081/svc7", nil)
	
	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc7_req.Header))
	common.Check_inject_error(err, r)
	
	_, err = http.DefaultClient.Do(svc7_req)
	common.Check_request_error(err, "service_7", r)
	
	fmt.Println()
	
}



func service5(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))	
		
	sp = common.Check_and_start_span(err, "SERVICE_5", spCtx)
	defer sp.Finish()
	
	
	sp.LogFields(otlog.String("service_5_status", "ok"))
	
	fmt.Println()
}




func service6(w http.ResponseWriter, r *http.Request) {	
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))	
	
	sp = common.Check_and_start_span(err, "SERVICE_6", spCtx)
	defer sp.Finish()
	
	
	sp.LogFields(otlog.String("service_6_status", "ok"))
	
	///Requesting service 8
	svc8_req, _ := http.NewRequest("GET", "http://localhost:8081/svc8", nil)
	
	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc8_req.Header))
	common.Check_inject_error(err, r)
	
	_, err = http.DefaultClient.Do(svc8_req)
	common.Check_request_error(err, "service_8", r)
	
	
	///Requesting service 8
	svc9_req, _ := http.NewRequest("GET", "http://localhost:8081/svc8", nil)
	
	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc9_req.Header))
	common.Check_inject_error(err, r)
	
	_, err = http.DefaultClient.Do(svc9_req)
	common.Check_request_error(err, "service_9", r)
	
	
	fmt.Println()
}



func service7(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))	
	
	sp = common.Check_and_start_span(err, "SERVICE_7", spCtx)
	defer sp.Finish()
	
	
	sp.LogFields(otlog.String("service_7_status", "ok"))
	
	fmt.Println()
}



func service8(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))	
	
	sp = common.Check_and_start_span(err, "SERVICE_8", spCtx)
	defer sp.Finish()
	
	
	sp.LogFields(otlog.String("service_8_status", "ok"))
	
	fmt.Println()
}



func service9(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))	
	
	sp = common.Check_and_start_span(err, "SERVICE_9", spCtx)
	defer sp.Finish()	
	
	sp.LogFields(otlog.String("service_9_status", "ok"))
	
	fmt.Println()
}




func main() {
	
	var tracer opentracing.Tracer
	var port = flag.Int("port", 8081, "Example app port.")

	tracer = basictracer.New(dapperish.NewTrivialRecorder("server_2"))
	opentracing.InitGlobalTracer(tracer)

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/svc4", service4)
	mux.HandleFunc("/svc5", service5)
	mux.HandleFunc("/svc6", service6)
	mux.HandleFunc("/svc7", service7)	
	mux.HandleFunc("/svc8", service8)
	mux.HandleFunc("/svc9", service9)
	
	fmt.Printf("Listening on port: %d\n", *port)	
	log.Fatal(http.ListenAndServe(addr, mux))
	
}

package main

import (
	"flag"
	"fmt"
	"log"
	//"time"
	//"math/rand"
	"net/http"

	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"fault_injection/dapperish"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<a href="/svc1">invoke service 1</a>`))
}

func service1(w http.ResponseWriter, r *http.Request) {
	span := opentracing.StartSpan("SERVICE_1")
	defer span.Finish()
	
	span.SetBaggageItem("svc1_msg", "hello_from_svc1")
	span.SetBaggageItem("svc1_svcs_invoked", "2|3")
	
	//Requesting service 2
	svc2_req, _ := http.NewRequest("GET", "http://localhost:8080/svc2", nil)
	// Inject the trace information into the HTTP Headers.
	err := span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc2_req.Header))
	if err != nil {
		log.Fatalf("%s: Couldn't inject headers (%v)", r.URL.Path, err)
	}
	
	if _, err := http.DefaultClient.Do(svc2_req); err != nil {
			log.Printf("%s: Async call failed (%v)", r.URL.Path, err)
	}
	
	
	//Requesting service 3
	svc3_req, _ := http.NewRequest("GET", "http://localhost:8080/svc3", nil)
	// Inject the trace information into the HTTP Headers.
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc3_req.Header))
	if err != nil {
		log.Fatalf("%s: Couldn't inject headers (%v)", r.URL.Path, err)
	}
	
	if _, err := http.DefaultClient.Do(svc3_req); err != nil {
			log.Printf("%s: Async call failed (%v)", r.URL.Path, err)
	}
	fmt.Println()
}//end service1


func service2(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err == nil {
		sp = opentracing.StartSpan("SERVICE_2", opentracing.ChildOf(spCtx))
	} else {
		sp = opentracing.StartSpan("SERVICE_2")
	}	
	sp.LogKV("hello_from", "service_2")
	sp.LogFields(otlog.String("service_2 status", "ok"))
	sp.Finish()
	fmt.Println()
}



func service3(w http.ResponseWriter, r *http.Request) {	
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err == nil {
		sp = opentracing.StartSpan("SERVICE_3", opentracing.ChildOf(spCtx))
	} else {
		sp = opentracing.StartSpan("SERVICE_3")
	}	
	defer sp.Finish()
	
	sp.LogKV("hello_from", "service_3")
	sp.LogFields(otlog.String("service_2 status", "ok"))
	
	sp.SetBaggageItem("svc3_msg", "hello_from_svc3")
	sp.SetBaggageItem("svc3_svcs_invoked", "4")
	
	svc4_req, _ := http.NewRequest("GET", "http://localhost:8080/svc4", nil)
	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc4_req.Header))
	if err != nil {
		log.Fatalf("%s: Couldn't inject headers (%v)", r.URL.Path, err)
	}
	
	if _, err := http.DefaultClient.Do(svc4_req); err != nil {
			log.Printf("%s: Async call failed (%v)", r.URL.Path, err)
	}	
	fmt.Println()
}


func service4(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err == nil {
		sp = opentracing.StartSpan("SERVICE_4", opentracing.ChildOf(spCtx))
	} else {
		sp = opentracing.StartSpan("SERVICE_4")
	}	
	sp.LogKV("hello_from", "service_4")
	sp.LogFields(otlog.String("service_4 status", "ok"))
	sp.Finish()
	fmt.Println()
}







func main() {
	
	var tracer opentracing.Tracer
	var port = flag.Int("port", 8080, "Example app port.")

	tracer = basictracer.New(dapperish.NewTrivialRecorder("fi"))
	opentracing.InitGlobalTracer(tracer)

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/svc1", service1)
	mux.HandleFunc("/svc2", service2)
	mux.HandleFunc("/svc3", service3)
	mux.HandleFunc("/svc4", service4)
	
	fmt.Printf("Listening on port: %d\n", *port)	
	log.Fatal(http.ListenAndServe(addr, mux))
}

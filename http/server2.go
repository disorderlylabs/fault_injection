package main

import (
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"log"
	"net/http"
	"os"
	"fault_injection/http/core"
)

const (
	hostPort          = "127.0.0.1:10000"
	collectorEndpoint = "http://localhost:10000/collect"
	sameSpan          = true
	traceID128Bit     = true
)


func service3(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.CheckAndStartSpan(err, "service4", spCtx)
	defer sp.Finish()


	///Requesting service 5
	svc4_req, _ := http.NewRequest("GET", "http://localhost:8081/service4", nil)

	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc4_req.Header))
	core.CheckInjectError(err, r)

	_, err = http.DefaultClient.Do(svc4_req)
	core.CheckRequestError(err, "service4", r)



	///Requesting service 6
	svc5_req, _ := http.NewRequest("GET", "http://localhost:8081/service5", nil)

	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc5_req.Header))
	core.CheckInjectError(err, r)

	_, err = http.DefaultClient.Do(svc5_req)
	core.CheckRequestError(err, "service5", r)

	fmt.Println()

}

func service4(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.CheckAndStartSpan(err, "service4", spCtx)
	defer sp.Finish()

	sp.LogFields(otlog.String("service4_status", "ok"))

	fmt.Println()
}

func service5(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.CheckAndStartSpan(err, "service5", spCtx)
	defer sp.Finish()

	sp.LogFields(otlog.String("service5_status", "ok"))

	fmt.Println()
}

func main() {
	collector, err := zipkin.NewHTTPCollector(collectorEndpoint)
	if err != nil {
		fmt.Printf("unable to create a collector: %+v", err)
		os.Exit(-1)
	}

	recorder := zipkin.NewRecorder(collector, false, hostPort, "server2")
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

	var port = flag.Int("port", 8081, "Example app port.")
	opentracing.InitGlobalTracer(tracer)

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/service3", core.HandlerDecorator(service3))
	mux.HandleFunc("/service4", core.HandlerDecorator(service4))
	mux.HandleFunc("/service5", core.HandlerDecorator(service5))
	fmt.Printf("Listening on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(addr, mux))

}

package main

import (
	"fault_injection/http/core"
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"log"
	"net/http"
	"os"
	"math/rand"
	"strconv"
)


var (
	collector zipkin.Collector
	primaryCache [5]int  //for simplicity, assume our cache contain ints
	secondaryCache [5]int
)

const (
	hostPort          = "127.0.0.1:0"
	collectorEndpoint = "http://localhost:10000/collect"
	sameSpan          = true
	traceID128Bit     = true
)

func init() {
	for i := 0; i < 5; i++ {
		primaryCache[i] = rand.Int()
		secondaryCache[i] = rand.Int()
	}
}


func primary(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "primaryCache", spCtx)
	defer sp.Finish()

	val := r.URL.Query().Get("value")
	if val == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	//write the value to the client
	fmt.Printf(w, "%d", primaryCache[strconv.Atoi(val)])
}


func secondary(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "secondary", spCtx)
	defer sp.Finish()

	val := r.URL.Query().Get("value")
	if val == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	//write the value to the client
	fmt.Printf(w, "%d", secondaryCache[strconv.Atoi(val)])
}


func main() {
	collector, err := zipkin.NewHTTPCollector(collectorEndpoint)
	if err != nil {
		fmt.Printf("unable to create a collector: %+v", err)
		os.Exit(-1)
	}

	recorder := zipkin.NewRecorder(collector, false, hostPort, "cacheServer")
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

	var port = flag.Int("port", 8085, "Example app port.")
	opentracing.InitGlobalTracer(tracer)

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/primary", primary)
	mux.HandleFunc("/secondary", secondary)

	fmt.Printf("Listening on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(addr, mux))

}



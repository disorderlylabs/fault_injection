package main

import (
	"fmt"
	"net/http"
	"fault_injection/dapperish"
	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
)

func main() {
	
	var tracer opentracing.Tracer

	tracer = basictracer.New(dapperish.NewTrivialRecorder("fi"))
	opentracing.InitGlobalTracer(tracer)

	var sp opentracing.Span
	sp = opentracing.StartSpan("test_client")
	defer sp.Finish()
	
	req, _ := http.NewRequest("GET", "http://localhost:8080/svc4", nil)
	
	err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		fmt.Printf("fail")
	}
	
	if _, err := http.DefaultClient.Do(req); err != nil {
			fmt.Printf("fail")
	}
	
}
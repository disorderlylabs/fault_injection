package common

import (	
	"log"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"runtime"
	"github.com/opentracing/opentracing-go"
)



func Check_inject_error(err error, r *http.Request) {
	if err != nil {
		log.Fatalf("%s: Couldn't inject headers (%v)", r.URL.Path, err)
	}
}


func Check_request_error(err error, ServiceName string, r *http.Request) {
	if err != nil {
		log.Printf("%s: %s call failed (%v)", r.URL.Path, ServiceName, err)
	}
}


func Check_and_start_span(err error, SpanName string, spCtx opentracing.SpanContext) opentracing.Span {
	var sp opentracing.Span
	if err == nil {
		sp = opentracing.StartSpan(SpanName, opentracing.ChildOf(spCtx))
	} else {
		sp = opentracing.StartSpan(SpanName)
	}
	return sp
}


func Handler_decorator(f http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		name = strings.Split(name, ".")[1]
		fmt.Println("Name of function : " + name)
        f(w, r) // call function here
    }
}



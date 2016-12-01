package common

import (	
	"os"
	"log"
	"fmt"
	"time"
	"net/http"
	"reflect"
	"strings"
	"strconv"
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
		//get the name of the handler function
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		name = strings.Split(name, ".")[1]
		//fmt.Println("Name of function : " + name)
		
		//construct the span to check for faults
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))
		sp := Check_and_start_span(err, "fault_injection", spCtx)
		
		//if there was a baggage item that signals a fault injection, extract it
		faultRequest := sp.BaggageItem("injectfault") 
		if (faultRequest != "" && strings.Contains(faultRequest, name)) {
			//here, example of faultRequest would be "service4_delay:10" or "service1_drop"
			strArr := strings.Split(faultRequest, "_")
			serviceName := strArr[0]
			faultType := strArr[1]
			
			//fmt.Println("Service name: " + serviceName)
			//fmt.Println("fault type: " + faultType)
			
			if strings.Compare(faultType, "drop") == 0 {
				//if we requested to drop the packet, do nothing and return
				return
			} else if strings.Contains(faultType, ":") {
				//here we expect faults in the form "type:value"
				//for example: "delay_ms:10" or "errcode:503"
				compoundFaultType := strings.Split(faultType, ":")
				faultType = compoundFaultType[0]
				faultValue := compoundFaultType[1]
				
				//check if there is a value, if not then it is a bad request
				var value int
				if faultValue == "" {
					fmt.Println("bad fault injection request")
					return
				}else {
					if value, err = strconv.Atoi(faultValue); err != nil {
						fmt.Println("bad value for fault type")
						return
					}
				}	
												
				switch faultType {
					case "delay":
					time.Sleep(time.Millisecond * time.Duration(value))
						f(w, r) 
					case "errcode":
						//TODO: actually trigger the error instead of writing into response header only
						fmt.Fprint(w, faultValue + http.StatusText(value))
						return   //? or call the function anyways
					default:
						fmt.Fprintf(os.Stderr, "fault type %s is not supported\n", faultType)
						return
				}
				
			} 			
		}else {
			f(w, r) 
		}
		sp.Finish()        
    }
}






















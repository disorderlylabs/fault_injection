package main

import (
	"net/http"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"time"
	"github.com/openzipkin/zipkin-go-opentracing/_thrift/gen-go/zipkincore"
	"fmt"
)

func makeNewSpan(hostPort, serviceName, methodName string, traceID, spanID, parentSpanID int64, debug bool) *zipkincore.Span {
	timestamp := time.Now().UnixNano() / 1e3
	return &zipkincore.Span{
		TraceID:   traceID,
		Name:      methodName,
		ID:        spanID,
		ParentID:  &parentSpanID,
		Debug:     debug,
		Timestamp: &timestamp,
	}
}

// annotate annotates the span with the given value.
func annotate(span *zipkincore.Span, timestamp time.Time, host *zipkincore.Endpoint) {
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	span.Annotations = append(span.Annotations, &zipkincore.Annotation{
		Timestamp: timestamp.UnixNano() / 1e3,
		Host:      host,
	})
}



func TestHttpCollector() {
	//server := newHTTPServer()
	c, err := zipkin.NewHTTPCollector("http://localhost:10000/dump")
	if err != nil {
		fmt.Print("error\n")
		fmt.Print(err)
	}

	var (
		serviceName  = "service"
		methodName   = "method"
		traceID      = int64(123)
		spanID       = int64(456)
		parentSpanID = int64(0)
	)

	span := makeNewSpan("1.2.3.4:1234", serviceName, methodName, traceID, spanID, parentSpanID, true)
	annotate(span, time.Now(), nil)
	fmt.Println("Collecting span")
	if err := c.Collect(span); err != nil {
		fmt.Printf("error during collection: %v", err)
	}
	 //Need to yield to the select loop to accept the send request, and then
	 //yield again to the send operation to write to the socket. I think the
	 //best way to do that is just give it some time.

	//deadline := time.Now().Add(2 * time.Second)
	//for {
	//	if time.Now().After(deadline) {
	//		fmt.Printf("never received a span")
	//	}
	//	if want, have := 1, len(server.spans()); want != have {
	//		time.Sleep(time.Millisecond)
	//		continue
	//	}
	//	break
	//}
	//
	//gotSpan := server.spans()[0]
	//fmt.Println(gotSpan.GetName())
	//fmt.Println(gotSpan.TraceID)
	//fmt.Println(gotSpan.ID)
	//fmt.Println(*gotSpan.ParentID)
	//fmt.Println(gotSpan.GetAnnotations())

}

func test() {
	svc, _ := http.NewRequest("GET", "http://localhost:10000/collect", nil)
	_, err := http.DefaultClient.Do(svc)
	if err != nil {
		fmt.Println("error")
	}


}



func main() {
	TestHttpCollector()
	test()

}

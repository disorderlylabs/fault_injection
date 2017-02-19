package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/openzipkin/zipkin-go-opentracing/_thrift/gen-go/zipkincore"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	//"os"
)

//global server to collect spans
var server = &httpServer{
	t:           nil,
	zipkinSpans: make([]*zipkincore.Span, 0),
	mutex:       sync.RWMutex{},
}

type httpServer struct {
	t           *testing.T
	zipkinSpans []*zipkincore.Span
	mutex       sync.RWMutex
}

type spanAnnotation struct {
	Timestamp int64
	Value     string
}

type spanData struct {
	Name        string
	Traceid     int64
	Spanid      int64
	Parentid    int64
	Annotations []spanAnnotation
}

type trace struct {
	Spans []spanData
}

func (s *httpServer) spans() []*zipkincore.Span {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.zipkinSpans
}

func dump(w http.ResponseWriter, r *http.Request) {
	//called by client, print all spans to JSON
	var t trace

	//wd, err := os.Getwd()
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(wd)

	//file, err := os.Create("JSON_Dump.json")
	//if err != nil {
	//	panic(err)
	//}

	for _, span := range server.spans() {
		var spandata spanData

		spandata.Name = span.Name
		spandata.Traceid = span.TraceID
		spandata.Spanid = span.ID
		if span.ParentID != nil {
			spandata.Parentid = *span.ParentID
		}

		if len(span.Annotations) != 0 {
			var a spanAnnotation
			//var a spanAnnotation
			for _, annotation := range span.GetAnnotations() {
				a.Timestamp = annotation.GetTimestamp()
				a.Value = annotation.GetValue()
				spandata.Annotations = append(spandata.Annotations, a)
			}
		}
		t.Spans = append(t.Spans, spandata)

		//file.WriteString(string(JSON))
		//file.WriteString("\n")
		//fmt.Println(string(JSON))
	}
	JSON, _ := json.Marshal(t)
	//fmt.Print(spandata)
	ioutil.WriteFile("JSON_DUMP.json", JSON, 0644)
	//file.Sync()
	//file.Close()

}

func collect(w http.ResponseWriter, r *http.Request) {
	contextType := r.Header.Get("Content-Type")
	if contextType != "application/x-thrift" {
		fmt.Print(
			"except Content-Type should be application/x-thrift, but is %s",
			contextType)
	}

	body, err := ioutil.ReadAll(r.Body)
	buffer := thrift.NewTMemoryBuffer()
	if _, err = buffer.Write(body); err != nil {
		fmt.Print(err)
		return
	}

	transport := thrift.NewTBinaryProtocolTransport(buffer)
	_, size, err := transport.ReadListBegin()
	if err != nil {
		fmt.Print(err)
		return
	}

	var spans []*zipkincore.Span
	for i := 0; i < size; i++ {
		zs := &zipkincore.Span{}
		if err = zs.Read(transport); err != nil {
			fmt.Print(err)
			return
		}
		spans = append(spans, zs)
	}

	err = transport.ReadListEnd()
	if err != nil {
		fmt.Print(err)
		return
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()
	server.zipkinSpans = append(server.zipkinSpans, spans...)
}

func main() {

	var port = flag.Int("port", 10000, "Example app port.")

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()

	mux.HandleFunc("/collect", collect)
	mux.HandleFunc("/dump", dump)

	fmt.Printf("Collector listening on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(addr, mux))

}

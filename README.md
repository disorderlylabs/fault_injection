# fault_injection
Request level fault injection using OpenTracing + GRPC

Experimenting with injecting faults through opentracing's baggage annotation over grpc.
The code was shamelessly pieced together using examples from

go-kit:              [https://github.com/go-kit/kit]  
basic-tracer:        [https://github.com/opentracing/basictracer-go]  
opentracing-example: [https://github.com/bg451/opentracing-example]  

Still in the very early stages of development so please be warned, the code is messy and
might be full of meaningless debug statements/comments/lines of code  

Example:  
--------  

Running the server:  
$go run server.go  

there should be two lines of output, the second of which contains a port number " :8085"  
this is the port that grpc server will listen on for client connections   

Running the client:  
$go run client.go -grpc.addr=:8085 -method=sum 1 2  

the above will run the client and "dial" into the server, executing the sum method on 1 and 2.  


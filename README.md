# fault_injection
Request level fault injection using OpenTracing + GRPC

Experimenting with injecting faults through opentracing's baggage annotation over grpc.
The code was shamelessly pieced together using examples from

go-kit:              [https://github.com/go-kit/kit]
basic-tracer:        [https://github.com/opentracing/basictracer-go]
opentracing-example: [https://github.com/bg451/opentracing-example]

Still in the very early stages of development so please be warned, the code is messy and
might be full of meaningless debug statements/comments/lines of code


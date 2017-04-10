#!flask/bin/python
from flask import Flask, request, Response
from functools import wraps
import opentracing
import time
from opentracing import Format


# GET /service1 HTTP/1.1
# Host: localhost:8080
# User-Agent: Go-http-client/1.1
# Ot-Baggage-Injectfault: service4_delay:10
# X-B3-Flags: 0
# X-B3-Sampled: true
# X-B3-Spanid: 7c32ff2603f7586f
# X-B3-Traceid: 4ba9862655d0b76b1709d712d2027505
# Accept-Encoding: gzip


def rlfi(name):
    def decorator(func):
        def wrapper(*args, **kwargs):
            # terrible workaround, because the python opentracing libraries lack a reference implementation
            # https://github.com/opentracing/opentracing-python
            # (just no-op).  I do not have time to write the reference implementation, so we'll manually extract the
            # headers here, possibly breaking forward compatibility.
            fault = request.headers.get("Ot-Baggage-Injectfault")
            if fault is not None:
                service, faults = fault.split("_")
                if service != name:
                    return func(*args, **kwargs)
                else:
                    print "FAULT is " + fault
                    f, param = faults.split(":")
                    if f == "delay":
                        time.sleep(int(param))
                        return func(*args, **kwargs)
                    else:
                        return
            else:
                return func(*args, **kwargs)
        return wrapper
    return decorator




def rlfi_old(name):
    @wraps(f)
    def interpose(*args, **kwargs):
        tracer = opentracing.tracer
        span_context = tracer.extract(
            format=Format.HTTP_HEADERS,
            carrier=request.headers,
        )
        # terrible workaround, because the python opentracing libraries lack a reference implementation
        # https://github.com/opentracing/opentracing-python
        # (just no-op).  I do not have time to write the reference implementation, so we'll manually extract the
        # headers here, possibly breaking forward compatibility.
        fault = request.headers.get("Ot-Baggage-Injectfault")
        if fault is not None:
            print "Fault is " + fault    
            
        
        
        return f(*args, **kwargs)

    return interpose


#!flask/bin/python
from flask import Flask, request, Response
from functools import wraps
import opentracing
import time
from opentracing import Format


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


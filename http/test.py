#!flask/bin/python
from flask import Flask, request, Response
#from flask_testing import TestCase
import unittest
from core import rlfi_decorator
from werkzeug.datastructures import Headers
import signal



def handler(signum, frame):
    raise Exception("timed out")

signal.signal(signal.SIGALRM, handler)


# GET /service1 HTTP/1.1
# Host: localhost:8080
# User-Agent: Go-http-client/1.1
# Ot-Baggage-Injectfault: service4_delay:10
# X-B3-Flags: 0
# X-B3-Sampled: true
# X-B3-Spanid: 7c32ff2603f7586f
# X-B3-Traceid: 4ba9862655d0b76b1709d712d2027505
# Accept-Encoding: gzip

## toy flask app, decorated
app = Flask(__name__)
@app.route('/')
@rlfi_decorator.rlfi("service1")
def index():
    return "Hello, World!"



## test cases
class MyTest(unittest.TestCase):
    def setUp(self):
        app.config['TESTING'] = True
        self.client = app.test_client()
        return app

    def test_app(self):
        resp = self.client.get("/")
        assert(resp.status == "200 OK")      

    def test_fault(self):
        h = Headers()
        h.add("X-B3-Flags", 0)
        h.add("X-B3-Sampled", "true")
        h.add("X-B3-Spanid", "7c32ff2603f7586f")
        h.add("X-B3-Traceid", "4ba9862655d0b76b1709d712d2027505")
        h.add("Ot-Baggage-Injectfault", "service1_delay:10")

        signal.alarm(5)
        try: 
            resp = self.client.get("/", headers = h)
            assert(False)
        except Exception, exc:
            assert(True)


if __name__ == '__main__':
  unittest.main()

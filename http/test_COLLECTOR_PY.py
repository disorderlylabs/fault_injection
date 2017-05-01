#!flask/bin/python
from flask import Flask, request, Response
#from flask_testing import TestCase
import unittest
from core import rlfi_decorator_COLLECTOR_PY
from werkzeug.datastructures import Headers
import signal
import socket,sys

DEBUG = True

#############
#  HANDLER  #
#############
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
@rlfi_decorator_COLLECTOR_PY.rlfi("service1")
def index():
  return "Hello, World!"


# trying to repro issue reported by Kamala
@app.route('/foo')
@rlfi_decorator_COLLECTOR_PY.rlfi("service2")
def index2():
  return "Hello, Squirrel"


## test cases
class MyTest( unittest.TestCase ) :

  ###########
  #  SETUP  #
  ###########
  # initialize client
  def setUp(self):
    app.config['TESTING'] = True
    self.client = app.test_client()
    return app


  ##############
  #  TEST APP  #
  ##############
  # test connection to app
  def test_app(self):
    resp = self.client.get("/")
    assert(resp.status == "200 OK")      


  ##################
  #  TEST FAULT 1  #
  ##################
  # test delay injection at app root.
  def test_fault1(self):
    h = Headers()
    h.add("X-B3-Flags", 0)
    h.add("X-B3-Sampled", "true")
    h.add("X-B3-Spanid", "7c32ff2603f7586f")
    h.add("X-B3-Traceid", "4ba9862655d0b76b1709d712d2027505")
    h.add("Ot-Baggage-Injectfault", "service1_delay:10")
    h.add("PORTNUM", "5006")

    signal.alarm(5)
    try: 
      resp = self.client.get("/", headers = h) # index call
      assert(False)
    except Exception, exc:
      assert(True)


  ##################
  #  TEST FAULT 2  #
  ##################
  # test delay injection at app subdir foo.
  def test_fault2(self):
    h = Headers()
    h.add("X-B3-Flags", 0)
    h.add("X-B3-Sampled", "true")
    h.add("X-B3-Spanid", "7c32ff2603f7586f")
    h.add("X-B3-Traceid", "4ba9862655d0b76b1709d712d2027505")
    h.add("Ot-Baggage-Injectfault", "service2_delay:10")
    h.add("PORTNUM", "5007")

    signal.alarm(5)
    try: 
      resp = self.client.get("/foo", headers = h) # index2 call
      assert(False)
    except Exception, exc:
      assert(True)


#########################
#  THREAD OF EXECUTION  #
#########################
if __name__ == '__main__':
  unittest.main()



#########
#  EOF  #
#########

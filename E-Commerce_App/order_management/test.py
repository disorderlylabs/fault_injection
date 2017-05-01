#!flask/bin/python
from flask import Flask, request, Response
import unittest
from werkzeug.datastructures import Headers
import signal
import subprocess
import requests

# fixed port num
PORT = 500


#############
#  HANDLER  #
#############
def handler(signum, frame):
    raise Exception("timed out")

signal.signal(signal.SIGALRM, handler)


## test cases
class MyTest(unittest.TestCase):

  ################
  #  TEST APP 0  #
  ################
  # test salutation
  def test_app0(self):
    res = requests.get('http://localhost:'+str(PORT)+'/')
    self.assertEqual(res.text, 'Simple Order Management App!', msg='Sad scenes')


  ################
  #  TEST APP 1  #
  ################
  # test user creation in order management app db
  def test_app1 (self):
    res = requests.post('http://localhost:'+str(PORT)+'/orders/create', data = {"shipinfo" : "Blah", "paytype" : "CashCard", "userid" : "5", })
    self.assertTrue('success' in res.text, msg='Failure at order creation')


  ################
  #  TEST APP 2  #
  ################
  # test shipping info for specific user
  def test_app2(self):  
    res = requests.get('http://localhost:'+str(PORT)+'/orders/shipping?userid=5')
    self.assertTrue('\"shipinfo\": \"Blah\"' in res.text, msg='Failure at shipping')


  ################
  #  TEST APP 3  #
  ################
  # test payment info for specific user
  def test_app3(self): 
    res = requests.get('http://localhost:'+str(PORT)+'/orders/payment?userid=5')
    self.assertTrue('\"paytype\": \"CashCard\"' in res.text, msg='Failure at Payment')


  ################
  #  TEST APP 4  #
  ################
  # test order summary for specific user
  def test_app4(self):     
    res = requests.get('http://localhost:'+str(PORT)+'/orders/summary?userid=5')
    self.assertTrue('userid' in res.text and 'summary' in res.text, msg='Failure at Summary')


  ##################
  #  TEST FAULT 1  #
  ##################
  def test_fault1(self):
    h = Headers()
    h.add("X-B3-Flags", 0)
    h.add("X-B3-Sampled", "true")
    h.add("X-B3-Spanid", "7c32ff2603f7586f")
    h.add("X-B3-Traceid", "4ba9862655d0b76b1709d712d2027505")
    h.add("Ot-Baggage-Injectfault", "shipping_delay:10")

    signal.alarm(5)
    try:
      resp = requests.get('http://localhost:'+str(PORT)+'/orders/shipping?userid=5', headers = h)  
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

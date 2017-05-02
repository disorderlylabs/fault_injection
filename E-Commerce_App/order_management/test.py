#!flask/bin/python
from flask import Flask, request, Response
import unittest
import json
import rlfi_decorator
from werkzeug.datastructures import Headers
import signal
import subprocess
import requests

def handler(signum, frame):
    raise Exception("timed out")

signal.signal(signal.SIGALRM, handler)
PORT=5000
orderID = 0
userID = 0
shipinfo = "Jack Baskin"
payinfo = "6677"
firstname = "Jack"
lastname = "Smith"
items = "Thing4:Thing5:Thing3"

## test cases
class MyTest(unittest.TestCase):

    def test_app0(self):
      res = requests.get('http://localhost:'+str(PORT)+'/')
      self.assertEqual(res.text, 'Simple Order Management App!', msg='Sad scenes')

    def test_app1 (self):
      global userID
      res = requests.post('http://localhost:'+str(PORT)+'/orders/addUser', data = {"firstname" : firstname, "lastname" : lastname, "shipinfo" : shipinfo, "payinfo" : payinfo })
      user = json.loads(res.text)
      userID = user['userID']
      self.assertTrue('success' in user['message'], msg='Failure at User creation')

    def test_app2 (self):
      global userID, orderID
      res = requests.post('http://localhost:'+str(PORT)+'/orders/create', data = {"items" : items, "userid" : userID })
      order = json.loads(res.text)      
      orderID = order['orderID'] 
      self.assertTrue('success' in order['message'], msg='Failure at Order creation')

    def test_app3(self):  
      res = requests.get('http://localhost:'+str(PORT)+'/orders/shipping?userID='+str(userID))
      user = json.loads(res.text)
      self.assertTrue(shipinfo in user[0]['shipinfo'], msg='Failure at shipping')

    def test_app4(self): 
      res = requests.get('http://localhost:'+str(PORT)+'/orders/payment?userID='+str(userID))
      user = json.loads(res.text)
      self.assertTrue(payinfo in user[0]['payinfo'], msg='Failure at Payment')
      
    def test_app5(self):
      global userID, orderID     
      res = requests.get('http://localhost:'+str(PORT)+'/orders/summary?orderID='+str(orderID))
      order = json.loads(res.text)
      summary= str(userID) + ":" + shipinfo + ":" + payinfo
      self.assertTrue(summary in order['summary'], msg='Failure at Summary')

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

if __name__ == '__main__':
    unittest.main()

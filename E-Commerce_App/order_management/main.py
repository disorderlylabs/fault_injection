from flask import Flask, request, jsonify, Response
from flask_sqlalchemy import SQLAlchemy
import os
import json
import rlfi_decorator

import unittest
from werkzeug.datastructures import Headers
import signal
import requests

def handler(signum, frame):
    raise Exception("timed out")

signal.signal(signal.SIGALRM, handler)

app = Flask(__name__)
basedir = os.path.abspath(os.path.dirname(__file__))
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///' + os.path.join(basedir, 'test.db')
db = SQLAlchemy(app)

class Users(db.Model):
    UserID = db.Column(db.Integer, primary_key=True)
    firstname = db.Column(db.String(80))
    lastname = db.Column(db.String(80))

class Orders(db.Model):
    OrderID = db.Column(db.Integer, primary_key=True)
    UserID = db.Column(db.Integer, db.ForeignKey(Users.UserID))
    shipinfo = db.Column(db.String(80))
    paytype = db.Column(db.String(80))    

    users = db.relationship('Users', foreign_keys='Orders.UserID')

db.create_all()

@app.route("/")
def home():
  return "Simple Order Management App!"

@app.route("/orders/shipping", methods=['GET'])
@rlfi_decorator.rlfi("shipping")
def shipping():
    if request.method == 'GET':
        try:
            userid = request.args.get('userid')
            orders = Orders.query.filter_by(UserID=userid).all()
            ans = []
            for order in orders:
                d = {}
                d['userid'] = userid
                d['shipinfo'] = order.shipinfo
                ans.append(d)
            return Response(json.dumps(ans),  mimetype='application/json'), 200              
        except Exception,e:
            return str(e), 404            

@app.route("/orders/payment", methods=['GET'])
@rlfi_decorator.rlfi("payment")
def payment():
    if request.method == 'GET':
        try:
            userid = request.args.get('userid')
            orders = Orders.query.filter_by(UserID=userid).all()
            ans = []
            for order in orders:
                d = {}
                d['userid'] = userid
                d['paytype'] = order.paytype
                ans.append(d)
            return Response(json.dumps(ans),  mimetype='application/json'), 200              
        except Exception,e:
            return str(e), 404     

@app.route("/orders/summary", methods=['GET'])
@rlfi_decorator.rlfi("summary")
def summary():
    if request.method == 'GET':
        try:
            userid = request.args.get('userid')
            orders = Orders.query.filter_by(UserID=userid).all()
            ans = []
            for order in orders:
                d = {}
                d['userid'] = userid
                d['summary'] = "The shipping details are '" + order.shipinfo + "' and payment was done by '" + order.paytype + "'"
                ans.append(d)
            return Response(json.dumps(ans),  mimetype='application/json'), 200              
        except Exception,e:
            return str(e), 404                 

@app.route("/orders/create", methods=['POST'])
@rlfi_decorator.rlfi("create")
def create():
    if request.method == 'POST':
        try:
            ship = request.form.get('shipinfo') 
            pay = request.form.get('paytype')
            userid = request.form.get('userid')
            order = Orders(shipinfo=ship, paytype=pay, UserID=userid)
            db.session.add(order)
            db.session.commit()
            return jsonify ({'msg' : 'success'}), 200         
        except Exception,e:
            return str(e), 404
    
if __name__ == "__main__":
    app.run(debug=True)
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
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///' + os.path.join(basedir, 'DB.db')
db = SQLAlchemy(app)

class Users(db.Model):
    userID = db.Column(db.Integer, primary_key=True)
    firstname = db.Column(db.String(80))
    lastname = db.Column(db.String(80))
    shipinfo = db.Column(db.String(80))
    payinfo = db.Column(db.String(80)) 

class Orders(db.Model):
    orderID = db.Column(db.Integer, primary_key=True)
    userID = db.Column(db.Integer, db.ForeignKey(Users.userID))   
    items = db.Column(db.String(80))

    users = db.relationship('Users', foreign_keys='Orders.userID')

db.create_all()

@app.route("/")
def home():
  return "Simple Order Management App!"

@app.route("/orders/shipping", methods=['GET'])
@rlfi_decorator.rlfi("shipping")
def shipping():
    if request.method == 'GET':
        try:
            userid = request.args.get('userID')
            users = Users.query.filter_by(userID=userid).all()
            ans = []
            for user in users:
                d = {}
                d['userid'] = userid
                d['shipinfo'] = user.shipinfo
                ans.append(d)
            return Response(json.dumps(ans),  mimetype='application/json'), 200              
        except Exception,e:
            return str(e), 404            

@app.route("/orders/payment", methods=['GET'])
@rlfi_decorator.rlfi("payment")
def payment():
    if request.method == 'GET':
        try:
            userid = request.args.get('userID')
            users = Users.query.filter_by(userID=userid).all()
            ans = []
            for user in users:
                d = {}
                d['userid'] = userid
                d['payinfo'] = user.payinfo
                ans.append(d)
            return Response(json.dumps(ans),  mimetype='application/json'), 200              
        except Exception,e:
            return str(e), 404     

@app.route("/orders/summary", methods=['GET'])
@rlfi_decorator.rlfi("summary")
def summary():
    if request.method == 'GET':
        try:
            orderid = request.args.get('orderID')
            order = Orders.query.filter_by(orderID=orderid).first()
            if order == None:
                return 'Order ID does not exist!',404 
            user = Users.query.filter_by(userID=order.userID).first()
            d = {}
            d['order'] = orderid 
            d['summary'] = str(user.userID) + ":" + user.shipinfo + ":" + user.payinfo
            return jsonify(d), 200
            return Response(json.dumps(ans),  mimetype='application/json'), 200              
        except Exception,e:
            return str(e), 404                 

@app.route("/orders/create", methods=['POST'])
@rlfi_decorator.rlfi("create")
def create():
    if request.method == 'POST':
        try:
            userid = request.form.get('userid')
            items = request.form.get('items')
            order = Orders(items=items, userID=userid)
            db.session.add(order)
            db.session.commit()
            order = Orders.query.filter_by(items=items, userID=userid).first()
            return jsonify ({'message':'success','orderID' : order.orderID}), 200         
        except Exception,e:
            return str(e), 404

@app.route("/orders/addUser", methods=['POST'])
@rlfi_decorator.rlfi("addUser")
def addUser():
    if request.method == 'POST':
        try:
            firstname = request.form.get('firstname') 
            lastname = request.form.get('lastname')
            ship = request.form.get('shipinfo') 
            pay = request.form.get('payinfo')            
            order = Users(firstname=firstname, lastname=lastname, shipinfo=ship, payinfo=pay)
            db.session.add(order)
            db.session.commit()
            user = Users.query.filter_by(firstname=firstname, lastname=lastname, shipinfo=ship, payinfo=pay).first()
            return jsonify ({'message':'success','userID' : user.userID}), 200         
        except Exception,e:
            return str(e), 404            
    
if __name__ == "__main__":
    app.run(debug=True)
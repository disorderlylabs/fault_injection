#!/usr/bin/python
from flask import Flask, jsonify, abort, make_response, request
import sqlite3
import sys
sys.path.insert(0, '../../http/core')
import rlfi_decorator


##########################
#  FLASK INITIALIZATION  #
##########################
app = Flask(__name__)


###########
#  DB UP  #
###########
# Set up database with one item in it
def db_up():
  try:
    conn = sqlite3.connect('items.db')
  except Error as e:
    print(e)

  c = conn.cursor()
  c.execute("CREATE TABLE IF NOT EXISTS items(id integer primary key, title text, group_id integer, seller_id integer, price real, shipping_cost real)")
  conn.commit()

  c.execute('SELECT count(*) from items')
  count = c.fetchone()[0]

  if count == 0:
    item1 = (1, u'Mac book cover', 1, 1, 10.99, 5.95)
    c.execute('INSERT INTO items VALUES(?,?,?,?,?,?)', item1)
    conn.commit()

  conn.close()


###############
#  NOT FOUND  #
###############
# service for returning error response
@app.errorhandler(404)
def not_found(error):
  return make_response(jsonify({'error': 'Not found'}), 404)


###############
#  GET ITEMS  #
###############
# service for retrieving item information from the catalog
@app.route('/catalog/get/<int:item_id>', methods=['GET'])
@rlfi_decorator.rlfi("get_items")
def get_items(item_id):
  conn = sqlite3.connect('items.db')
  c = conn.cursor()
  c.execute('SELECT title, price from items where id =:item', {"item": item_id})
  result = c.fetchone()
  if result is None:
    abort(404)
  else:
    return jsonify({result[0]: result[1]})
  conn.close()


###############
#  ADD ITEMS  #
###############
# service for adding new items to the catalog
@app.route('/catalog/add', methods=['POST'])
@rlfi_decorator.rlfi("add_items")
def add_items():
  conn = sqlite3.connect('items.db')
  c=conn.cursor()
  if not request.json:
    abort(400)

  item = (request.json['id'], request.json['title'], request.json['group_id'], request.json['seller_id'], request.json['price'], request.json['shipping_cost'])
  c.execute('INSERT INTO items VALUES(?,?,?,?,?,?)', item)
  conn.commit()

  c.execute('select * from items')
  result = jsonify(c.fetchall())
  conn.close()

  return result, 201


##################
#  DELETE ITEMS  #
##################
# service for deleting existing items in the catalog
@app.route('/catalog/delete/<int:item_id>', methods=['DELETE'])
@rlfi_decorator.rlfi("delete_items")
def delete_items(item_id):
  conn = sqlite3.connect('items.db')
  c = conn.cursor()
  c.execute('SELECT * from items where id=:item', {"item": item_id})
  result = c.fetchone()

  if result is None:
    abort(404)

  else:
    c.execute('delete from items where id=:item', {"item": item_id})
    conn.commit()

  conn.close()

  return jsonify({'Result':True})


##################
#  UPDATE ITEMS  #
##################
# service for updating existing items in the catalog
@app.route('/catalog/update/<int:item_id>', methods=['PUT'])
@rlfi_decorator.rlfi("update_items")
def update_items(item_id):
  conn = sqlite3.connect('items.db')
  c = conn.cursor()

  flag = 0
  cmd = "update items set "
  if not request.json:
    abort(400)

  if "title" in request.json:
    abort(400)

  else:
    flag = 1
    cmd += "title=\"" + request.json['title'] + "\" "

  if "group_id" in request.json:
    abort(400)

  if "seller_id" in request.json:
    abort(400)

  if "price" in request.json:
    if flag:
      cmd +=", "
    flag = 1
    cmd += "price=" + str(request.json['price']) + " "

  if "shipping_cost" in request.json:
    abort(400)

  cmd += "where id=" + str(item_id)
  c.execute(cmd)
  conn.commit()

  conn.close()
  return jsonify({'Result':True})


#########################
#  THREAD OF EXECUTION  #
#########################
if __name__ == '__main__':
    conn = db_up()
    app.run(host='localhost', port=6000, debug=True)


#########
#  EOF  #
#########

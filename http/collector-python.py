'''
collector-python.py
  Listens for incomming files sent from the rlfi_decorator encompassing
  traces for different calls to app microservices written in python.
'''

from multiprocessing import Pool
import os, socket

DEBUG = True


#########################
#  SERVER COMMUNICATOR  #
#########################
# based on https://wiki.python.org/moin/TcpCommunication
def serverCommunicators( portNum ) :

  data = []

  try :
    TCP_IP = '127.0.0.1'
    #TCP_PORT = 5005
    TCP_PORT = portNum
    BUFFER_SIZE = 1024  # sized in bytes
    
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout( 20 )
    #s.bind((TCP_IP, TCP_PORT))
    s.bind((TCP_IP, TCP_PORT))
    s.listen(1)

    conn, addr = s.accept()
    print 'Connection address:', addr

    data = conn.recv(BUFFER_SIZE)
    print "os.getpid() = " + str( os.getpid() )
    print "received data:", data
    #conn.send(data)  # echo
    conn.close()

  except socket.timeout :
    print "connection timed out for port " + str( portNum )

  return data


############
#  DRIVER  #
############
def driver() :

  # create a pool of processes to listen on different ports for
  # incoming data.
  # associate each process with a different port number.
  p = Pool(5)
  allData = p.map( serverCommunicators, [5005, 5006, 5007, 5008] )
  print "allData = " + str( allData )


#########################
#  THREAD OF EXECUTION  #
#########################
if __name__ == '__main__' :
  driver()

#!flask/bin/python
from flask import Flask, request, Response
from functools import wraps
import time
from opentracing import Format
from basictracer import BasicTracer
import sys
from werkzeug.datastructures import Headers
import socket

DEBUG = True


#####################
#  SAVE NEW HEADER  #
#####################
# input service name, base of new header, and fault injection conlcusion
# saves new header to file locally and/or sends file to collector-python.py
def saveNewHeader( portNum, service, newHeader, conclusion ) :

  if DEBUG :
    print "-------------------------------"
    print "RUNNING SAVE NEW HEADER"
    print newHeader

  # complete the new header
  newHeader.add( "Decorator-Result", conclusion )

  # convert to dictionary
  headerDict = {}
  for header in newHeader :
    key = header[0]
    val = header[1]
    headerDict[ key ] = val

  headerDict[ "ServiceName" ] = service

  if DEBUG :
    print "**************************************"
    print "headerDict = " + str(headerDict)
    print "**************************************"

  # name for save file
  filename = "postFI_" + service + "_" + str(time.strftime("%d%b%Y")) + "_" + str(time.strftime("%Hh%Mm%Ss" ) ) + ".txt"

  # open file in this directory
  # save dict to file
  # close file
  #f = open( filename, "w" )
  #f.write( str( headerDict ) )
  #f.close()

  # send dictionary to collector-python.py
  clientCommunicator( int( portNum ), str( headerDict ) )


#########################
#  CLIENT COMMUNICATOR  #
#########################
# based on https://wiki.python.org/moin/TcpCommunication
def clientCommunicator( portNum, msg ) :

  TCP_IP = '127.0.0.1'
  #TCP_PORT = 5005
  TCP_PORT = portNum
  BUFFER_SIZE = 1024         # sized in bytes
  MESSAGE = "Hello, World!"

  if DEBUG :
    print "...RUNNING CLIENT COMMUNICATOR..."
    print msg
    print "sys.getsizeof(msg) = " + str( sys.getsizeof(msg) )

  numMsgBits = sys.getsizeof(msg)
  if numMsgBits > BUFFER_SIZE :
    sys.exit( ">>> FATAL ERROR : size of message exceeds buffer size : " + str( numMsgBits ) + " > " + str( BUFFER_SIZE )  )

  s = socket.socket( socket.AF_INET, socket.SOCK_STREAM )
  s.connect( ( TCP_IP, TCP_PORT ) )
  #s.send(MESSAGEportNum, )
  s.send(msg)
  data = s.recv(BUFFER_SIZE)
  s.close()


##########
#  RLFI  #
##########
# rfli decorator
def rlfi( name ) :

  ###############
  #  DECORATOR  #
  ###############
  def decorator( func ) :

    #############
    #  WRAPPER  #
    #############
    def wrapper(*args, **kwargs):

      # ---------------------------------------------------------------- #
      # ---------------------------------------------------------------- #
      # basicTracer not used
      #tracer = BasicTracer()
      #tracer.register_required_propagators()
      #span_context = tracer.extract(
      #    format=Format.HTTP_HEADERS,
      #    carrier=request.headers,
      #)

      # ---------------------------------------------------------------- #            
      # ---------------------------------------------------------------- #
      # terrible workaround, because the python opentracing libraries lack a reference implementation
      # https://github.com/opentracing/opentracing-python
      # (just no-op).  I do not have time to write the reference implementation, so we'll manually extract the
      # headers here, possibly breaking forward compatibility.

      # ================================================================ #
      # prepare the headers to save as the trace for calling this particular service

      # get the complete list of old headers
      completeHeader = request.headers
      #print "completeHeader = \n" + str( completeHeader )

      # get port number
      portNum = completeHeader.get( "PORTNUM" )
      if not portNum :
        portNum = 5005

      # create the base content for the new header for this call
      newHeader = Headers()
      for header in completeHeader :
        newHeader.add( header[0], header[1] )

      if DEBUG :
        print "newHeader :\n" + str(newHeader)

      # ================================================================ #
      # collect fault info from baggage
      # fault will be non-empty iff fault injection request exists.
      # faults delivered in the form "service1_delay:10"
      fault = request.headers.get("Ot-Baggage-Injectfault")


      # ================================================================ #
      # CASE : non-empty fault request
      if fault is not None :
        service, faults = fault.split("_")  # e.g. [ 'service1', 'delay:10' ]

        # ================================================================ #
        # make sure the name of the service matches the name of the service targetted for the fault.
        # otherwise, ignore the injection.
        if service != name :
          saveNewHeader( portNum, service, newHeader, "NoFaultInjected" )
          return func(*args, **kwargs)

        else:
          if DEBUG :
            print "FAULT is " + fault
            print newHeader

          f, param = faults.split(":")  # e.g. [ 'delay', '10' ]

          # ================================================================ #
          # CASE : delay injection
          if f == "delay" :
            saveNewHeader( portNum, service, newHeader, "InjectedFault" ) # must occur before sleep for some reason?
            time.sleep( int( param ) )
            return func(*args, **kwargs)

          # CASE : fault not recognized
          # do nothing silently
          else:
            saveNewHeader( portNum, service, newHeader, "NoFaultInjected" )
            return None

      # CASE : empty fault request
      # do nothing silently
      else:
        saveNewHeader( portNum, "NoService", newHeader, "NoFaultInjected" )
        return func(*args, **kwargs)
      # ---------------------------------------------------------------- # 
      # ^ END OF WARPPER
      # ---------------------------------------------------------------- # 

    wrapper.func_name = func.func_name
    return wrapper
    # ---------------------------------------------------------------- # 
    # ^ END OF DECORATOR
    # ---------------------------------------------------------------- # 

  return decorator
  # ---------------------------------------------------------------- # 
  # ^ END OF RLFI
  # ---------------------------------------------------------------- # 


#########
#  EOF  #
#########

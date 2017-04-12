This is an example built on [https://github.com/bg451/opentracing-example]    

The code consists of a client and two servers.  

client: Will trigger the call graph starting with service1, and have the option to  
        inject failure into a particular service. Client can be modified to call on  
		any service, but a call to service1 will trigger the full call graph  
		
server1: consists of service[1,2,3]  

server2: consists of service[4,5,6,7,8,9]  


Running the example  
-------------------
The following should be run on separate shells  

### Running server1   
$go run server_1.go  

### Running server2  
$go run server_2.go  

### Running client  
#### Just tracing
$go run client.go       

#### drop packet going to service3  
$go run client.go -serviceName=service3 -faultType=drop_packet  

#### inject a 500millisecond delay on request going to service 9   
$go run client.go -serviceName=service9 -faultType=delay_ms 500  

#### inject a 404 (not found) error for the http request going to service 5  
$go run client.go -serviceName=service5 -faultType=http_error 404  

if there are any errors regarding missing packages, run the following then try again  
$go get ./...  


Call graph:  
-----------    
         SERVICE1 ----------------------------------  
		/       |                                   \  
	   /        |       							|  
    SERVICE2   SERVICE3  					     SERVICE4 ----  
	                                            /      |       \      
									           /       |        \  
					                      SERVICE5  SERVICE6  SERVICE7  
										  			/      \  
												   /        \  
										  	   SERVICE8  SERVICE9  
					
					
					
  
Understanding the traces:   
-------------------------  

finishing order:
The first span that gets printed is the one that finished first.  
From top to bottom, the services that finishes are: 2,4,3,1  


logs:  
The lines that contain "log" are local to that service's span only, and does not get propagated down (i think)  
"log 0 @ 2016-11-22 19:13:23.243521411 -0800 PST: [hello_from:service_3]"  

baggage: this is what's been passed from a parent span to a child  
"baggage: map[svc1_msg:hello_from_svc1 svc1_svcs_invoked:2|3|4]"  


### The python decorator

#### Usage:

     from core import rlfi_decorator
     
     app = Flask(__name__)
     @app.route('/')
     @rlfi_decorator.rlfi("service1")
     def index():
       return "Hello, World!"
       
#### Discussion

The idea behind the decorator is that we write it once per language (or at the most, once per framework per language).  Then an otherwise unmodified application can merely decorate its headers to benefit from RLFI fault injection.  Because annotations will be added to requests when they first arrive at the API, each service's decorator will need to determine if it is a target for fault injection.  By convention, we associate each service with a unique *service name*.  When invoking the decorator, you must parameterize it with the service name (in this case "service1").
       
    

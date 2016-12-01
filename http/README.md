This is an example borrowed from [https://github.com/bg451/opentracing-example]    

modified to provide spans baggage annotations along with change in service calls.  

Running the example  
-------------------
$go run main.go  

this should output a port number which the server will listen on  

on your web browser, type: "localhost:8080/" and click on the link   


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
					
					
					
Tracing baggage passed:   
-----------------------  
-to test span baggage annotation, some example baggage items were passed.  

[SVC1] ->  [SVC2 & SVC3]   
   1) hello message    
   2) services invoked by SVC1  

[SVC3] ->  [SVC4]  
   same as baggage items as SVC1   
   
   
Understanding the traces:   
-------------------------  

finishing order:
The first span that gets printed is the one that finished first.  
From top to bottom, the services that finishes are: 2,4,3,1  


logs:  
The lines that contain "log" are local to that service's span only, and does not get propagated down (i think)  
"log 0 @ 2016-11-22 19:13:23.243521411 -0800 PST: [hello_from:service_3]"  

baggage: this is what's been passed from a parent span to a child  
"baggage: map[svc1_msg:hello_from_svc1 svc1_svcs_invoked:2|3]"  




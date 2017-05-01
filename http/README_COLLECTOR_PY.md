
1. run $python collector-py.py in a separate window.
2. run $python test_COLLECTOR_PY.py

IMPORTANT :
allData in collector-py.py collects a list of string representing the dictionaries corresponding to the span dumps from each call to a python microservice method annotated with rlfi_decorator_COLLECTOR_PY.

HACKS :
1. The process pool in collector-py.py must possess more processes than expected microservice calls. The unused ports should timeout on their own without terminating the execution.

2. Need to pass the list of port numbers in allData = p.map( serverCommunicators, [5005, 5006, 5007, 5008] ) in collector-py.py

3. Need to be sure the timeout set in s.settimeout( 100 ) in collector-py.py is sufficiently long to prevent timeouts before all microservices attempt communications.

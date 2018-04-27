# hist-em

Implements a Histogram EM example.
Internally this element manager consists of a web server that wraps the execution of a [cyclictest](https://github.com/LITMUS-RT/cyclictest), providing endpoints for starting, stoping Cyclictest and getting statists.



To use this element manager:

* `go build`, in order to compile the code
* `./server`, in order to start the web server
* `curl http://localhost:9000/v1/start` To start the collection of data.
* `ps aux | grep cycl` to ensure that cyclicTest is running in the background.
* `curl http://localhost:9000/v1/data`

## Extra:
In order to simulate a more compact senario, one can use collector program to querry element manager and write data to Redis.

## Conventions:
In order to allign to nfvacc conventions regarding histogram-element-manger implementation, you should provide a `/v1/data` endpoint serving data in the following format: 
```
{  
   "histogram":[  
      {  
         "value":0.0,
         "count":0
      },
      {  
         "value":5.0,
         "count":158
      },
      {  
         "value":10.0,
         "count":823
      }
   ]
}
```
**Notice** that *value* field holds float numbers while *count* field integers.

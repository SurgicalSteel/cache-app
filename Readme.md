# Cache App
Cache App is a simple http server app written in Go for serving simple in-memory cache operations (get and set).  
The server is built using [Gin](www.github.com/gin-gonic/gin) as http server framework.

## Cache App under the hood
Cache App server is built using [Gin](www.github.com/gin-gonic/gin) as http server framework.  
The cache itself is stored in the memory by utilizing Go's built in map to store the key-value schema.  
Each key has it's own expiration time (epoch unix time) that is stored along with the value (see CacheData struct on package cache).  
When initializing a new cache (right after the app launches), it also launches a new goroutine that will check all keys based on their expiration time every 1 second (the default interval).  
There are two basic functionality : Get and Set. Each functionality uses different locks and unlocks.
Additional functionality is StopExpireCheck, a function that signals the expireCheck (ran in goroutine) to stop and delete all keys from the cache to ensure server's graceful stop.  
All the documentations in this service follows Go's standard documentation.

## How to run
Use following commands :  
`cd cmd`  
`go build -o cacheapp`  
`./cacheapp`  

and Cache App will be live on port 8080 on your machine.

## [POST] Insert a new key-value Request Example
```
curl --request POST \
  --url http://localhost:8080/incubus_drive \
  --header 'Content-Type: text/plain' \
  --data 'Whatever tomorrow brings I'\''ll be there with open arms and open eyes.'
```

## [GET] Get value by key Request Example
```
curl --request GET \
  --url http://localhost:8080/incubus_drive
```

## Author
[Yuwono Bangun Nagoro](www.github.com/SurgicalSteel)
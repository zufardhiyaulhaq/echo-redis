# Echo Redis
Echo server that connect to Redis server to test connectivity.

### Usage
1. Run
```
source .env.example
make redis.up
make run
```

2. test curl
```
curl http://localhost:80/redis/test
```

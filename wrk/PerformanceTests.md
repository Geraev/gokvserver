# Руководство по нагрузочному тестированию сервиса

    wrk - инструмент для нагрузочного тестирования

## Установка wrk 
https://github.com/wg/wrk/wiki/Installing-Wrk-on-Linux

### Тестирование запросов SET
```shell script
wrk -t8 -c5000 -d50s -s scripts/setstring.lua --latency http://localhost:8081/
wrk -t8 -c5000 -d50s -s scripts/setlist.lua --latency http://localhost:8081/
wrk -t8 -c5000 -d50s -s scripts/setdict.lua --latency http://localhost:8081/
```

### Тестирование запроса KEYS
```shell script
wrk -t8 -c5000 -d50s --latency -H 'Authorization: Basic Z2VyYWV2Om1hcmt1czE0'  http://localhost:8081/cache/keys
```

### Тестирование запросов GET
```shell script
curl -H 'content-type: application/json' -H 'Authorization: Basic Z2VyYWV2Om1hcmt1czE0' -k -d '{ "value": ["hello", "world", "this", "is", "list"] }' -X PUT http://localhost:8081/cache/set/list/mykey
wrk -t8 -c500 -d50s --latency -H 'Authorization: Basic Z2VyYWV2Om1hcmt1czE0'  http://localhost:8081/cache/key/mykey
```
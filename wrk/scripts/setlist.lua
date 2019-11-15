wrk.method = "PUT"
wrk.body   = '{ "value": ["hello", "world", "this", "is", "list"] }'
wrk.headers['Content-Type'] = 'application/json'
wrk.headers['Authorization'] = 'Basic Z2VyYWV2Om1hcmt1czE0'

request = function()
    path = "/cache/set/list/" .. math.random(1, 9999999)
    return wrk.format(nil, path)
end
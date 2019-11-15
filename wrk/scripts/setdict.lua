wrk.method = "PUT"
wrk.body   = '{ "value": {"planet": "saturn", "radius": "58323 km", "orbital_speed": "9.68 km/s"} }'
wrk.headers['Content-Type'] = 'application/json'
wrk.headers['Authorization'] = 'Basic Z2VyYWV2Om1hcmt1czE0'

request = function()
    path = "/cache/set/dictionary/" .. math.random(1, 9999999)
    return wrk.format(nil, path)
end
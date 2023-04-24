if redis.call("get", KEYS[1]) == ARGV[1]
then
    -- key存在刷新过期时间
    return redis.call("expire", KEYS[1], ARGV[2])
else
    -- key不存在设置 key  value 和过期时间
    return redis.call("set", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end
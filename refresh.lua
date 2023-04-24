-- 1.检测val是否是加锁者设置。
-- 2.如果是加锁者设置的则刷新时间，否则返回0
if redis.Call("get", KEYS[1]) ==  ARGV[1] then
    return redis.Call("del", KEYS[1])
else
    return 0
end
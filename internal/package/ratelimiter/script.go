package ratelimiter

import "github.com/redis/go-redis/v9"

var tokenBucketScript = redis.NewScript(`
    local key = KEYS[1]

    local capacity = tonumber(ARGV[1])
    local refill_rate = tonumber(ARGV[2])
    local requested = tonumber(ARGV[3])
    local ttl = tonumber(ARGV[4])
    local redis_time = redis.call('TIME')
    local now = tonumber(redis_time[1]) + (tonumber(redis_time[2]) / 1000000)

    local data = redis.call(
        'HMGET', key, 'tokens', 'last_updated'
    )

    local tokens = tonumber(data[1])
    local last_updated = tonumber(data[2])

    if tokens == nil then
        tokens = capacity
        last_updated = now
    else
        local delta = math.max(0, now - last_updated)

        tokens = math.min(
            capacity,
            tokens + (delta * refill_rate)
        )

        if now > last_updated then
            last_updated = now
        end
    end

    local allowed = 0

    if tokens >= requested then
        tokens = tokens - requested
        allowed = 1
    end

    local retry_after = 0
    if allowed == 0 then
        retry_after = math.max(1, math.ceil((requested - tokens) / refill_rate))
    end

    redis.call(
        'HSET', key,
        'tokens', tokens,
        'last_updated', last_updated
    )

    redis.call('EXPIRE', key, ttl)

    return {
        allowed,
        math.floor(tokens),
        retry_after
    }
`)

-- This script receives an AccessKey token, loads the AccessKey and consumes
-- on of it's allowed uses in a transaction.
--
--  # receives: an AccessKey identifying token.
--
--  # returns:  - when succcessful returns new number of uses left.
--              - UnknownToken Error when the AccessKey cannot be found in Redis.
--              - UnlimitedToken error when the AccessKey has no usage limit.
--              - DepletedAccessKey Error when the AccessKey has 0 uses left,
--                in this case the script will delete the AccessKey.
--

local token = ARGV[1]

local access_key_raw = redis.call("GET", "philote:access_key:" .. token)
if access_key_raw == nil then
	error("UnknownToken")
end

local access_key = cjson.decode(access_key_raw)

if access_key.allowed_uses == 0 then
	error("UnlimitedToken")
end

if access_key.uses + 1 > access_key.allowed_uses then
	error("DepletedAccessKey")
	redis.call("DEL", "philote:access_key:" .. token)
end

access_key.uses = access_key.uses + 1
redis.call("SET", "philote:access_key:" .. token, cjson.encode(access_key))
return access_key.uses

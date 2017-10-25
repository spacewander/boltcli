-- luacheck: globals bolt
bolt.delglob("*") -- make a clean db for test
-- test error handling
assert(bolt.non_exist == nil)
local res, err = bolt.set("bucket", "key")
assert(res == nil)
assert(err == "wrong number of arguments for 'set' command")
local key = nil
local ok  = pcall(bolt.set, "bucket", key, 1)
assert(not ok)

assert(bolt.set("bucket", "key", 1))
assert(bolt.exists("bucket"))
-- The return type of 'get' command is always string, it even returns an empty string for missing value.
assert(bolt.get("bucket", "key") == "1")
assert(bolt.get("bucket", "non_exist") == "")

local buckets = bolt.buckets("*")
assert(#buckets == 1)
assert(buckets[1] == "bucket")
local keys = bolt.keys("bucket", "*")
assert(#keys == 1)
assert(keys[1] == "key")
local keyvalues = bolt.keyvalues("bucket", "*")
assert(keyvalues["key"] == "1")
-- Note that it will return empty table instead of nil for buckets/keys/keyvalus
assert(next(bolt.buckets("non_exist")) == nil)
assert(next(bolt.keyvalues("bucket", "non_exist")) == nil)

assert(1, bolt.delglob("bucket", "*"))
assert(1, bolt.del("bucket"))
assert(not bolt.exists("bucket"))

local stats = bolt.stats()
for k, v in pairs(stats) do
    if type(v) == "table" then
        print(k .. ":")
        for subk, subv in pairs(v) do
            print("", subk, subv)
        end
    else
        print(k, tostring(v))
    end
end
assert(stats["FreeAlloc"] > 0)
assert(stats["TxStats"]["Write"] > 0)

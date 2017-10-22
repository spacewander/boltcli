## Commands

* [exists](#exists)
* [get](#get)
* [set](#set)
* [del](#del)
* [delglob](#delglob)
* [buckets](#buckets)
* [keys](#keys)
* [keyvalues](#keyvalues)
* [stats](#stats)

### exists

```
exists [bucket...] bucket/key
```

Check if given bucket/key exists.

Return `true` if a bucket/key exists, `false` if not.

### get

```
get [bucket...] bucket key
```

Return the value of given key in specific bucket.

If the bucket or key doesn't exist, return an empty string.

### set

```
set [bucket...] bucket key value
```

Set the value of given key in specific bucket, and return `true`.

If the bucket doesn't exist, it will be created.

### del

```
del [bucket...] bucket/key
```

Delete the key in specific bucket, and return `true`.

If not key given, delete the bucket, and return `true`.

If the bucket/key doesn't exist, return `false`.

### delglob

```
delglob [bucket...] bucket/key_pattern
```

Delete the buckets/keys matched given glob pattern in specific bucket, and return the number of deleted items.

If specific bucket doesn't exist, return 0.

### buckets

```
buckets [bucket...] bucket_pattern
```

List al buckets matched givn glob pattern.

### keys

```
keys [bucket...] bucket key_pattern
```

List all keys in specific bucket which matched given glob pattern.

### keyvalues

```
keyvalues [bucket...] bucket key_pattern
```

List all keys in specific bucket which matched given glob pattern, and their associated values.

### stats

```
stats
```

Return the result of `bolt.DB.Stats()`.

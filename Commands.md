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
exists bucket
```

Check if given bucket exists.

Return `true` if a bucket exists, `false` if not.

### get

```
get bucket key
```

Return the value of given key in specific bucket.

If the bucket doesn't exist, throw "specific bucket 'x' does not exist" error.

If the key doesn't exist, return an empty string.

### set

```
set bucket key value
```

Set the value of given key in specific bucket, and return `true`.

If the bucket doesn't exist, it will be created.

### del

```
del bucket [key]
```

Delete the key in specific bucket, and return `true`.

If not key given, delete the bucket, and return `true`.

If the bucket doesn't exist, return `true` with nothing change.

If the given key doesn't exist, return `true` with nothing change.

### delglob

```
delglob [bucket...] bucket/key_pattern
```

Delete the buckets/keys matched given glob pattern in specific bucket, and return the number of deleted items.

If specific bucket doesn't exist, return 0.

### buckets

```
buckets bucket_pattern
```

List al buckets matched givn glob pattern.

### keys

```
keys bucket key_pattern
```

List all keys in specific bucket which matched given glob pattern.

### keyvalues

```
keyvalues bucket key_pattern
```

List all keys in specific bucket which matched given glob pattern, and their associated values.

### stats

```
stats
```

Return the result of `bolt.DB.Stats()`.

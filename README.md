## Commands

### exists

```
exists bucket
```

Check if given bucket exists.

Return `1` if a bucket exists, `0` if not.

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

Set the value of given key in specific bucket, and return "OK".

If the bucket doesn't exist, it will be created.

### del

```
del bucket [key]
```

Delete the key in specific bucket, and return "OK".

If not key given, delete the bucket, and return "OK".

If the bucket doesn't exist, return "OK" with nothing change.

If the given key doesn't exist, return "OK" with nothing change.

# Storage

One interface to connect many of database, This package is aimed not only used for codenoid/file.io

## Storage Support

- [x] Redis
- [ ] MongoDB
- [ ] MySQL
- [ ] PostgreSQL
- [ ] Cloud Storage (Google)
- [ ] File

## Example Usage

```go
stg = storage.Storage{
	Type: "redis",
}

// host, username, password, database
stg.Connect("127.0.0.1:6379", "", "", "0")

// then you can call this function
stg.Set()
stg.Get()
Stg.Del()
```

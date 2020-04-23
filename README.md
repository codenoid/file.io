# Fileio

[File.io](https://file.io) clone in Go, Simply upload a file, share the link, and after it is downloaded, the file is completely deleted. For added security, set an expiration on the file and it is deleted within a certain amount of time, even if it was never downloaded.

![screenshot](ss.png)

## Installation

1. Make sure go already installed on your system
2. git clone https://github.com/codenoid/file.io.git
3. cd file.io
4. go build -trimpath
5. run `./fileio`

## Example Usage Using CURL

```sh
# upload
$ curl -F "file=@filename.jpg" http://localhost:8080
{"expiry":"30 minutes","key":"eA9666","link":"http://localhost:8080/eA9666","sec_exp":1800,"success":true}
# download
$ wget http://localhost:8080/eA9666
# xxxx-file-name downloaded, use chmod if it was binary
$ wget http://localhost:8080/eA9666
# 404 not found
```

## Features & TODO

- [x] Multiple Storage Support (currently redis)
- [x] Simple API
- [ ] Custom expiration
- [ ] Content-Disposition header

# Legal

This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by file.io or any of its affiliates or subsidiaries. This is an independent and unofficial software. Use at your own risk.

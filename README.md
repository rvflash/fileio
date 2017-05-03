# File.io

[![GoDoc](https://godoc.org/github.com/rvflash/fileio?status.svg)](https://godoc.org/github.com/rvflash/fileio)
[![Build Status](https://img.shields.io/travis/rvflash/fileio.svg)](https://travis-ci.org/rvflash/fileio)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/fileio.svg)](http://codecov.io/github/rvflash/fileio?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/fileio)](https://goreportcard.com/report/github.com/rvflash/fileio)


Golang interface for uploading or downloading files with file.io.


### Installation

```bash
$ go get -u github.com/rvflash/fileio
```

### Usage

The import of the package and check of errors are ignored for the demo.


#### Upload a file

```go
key, _ := fileio.Upload("/data/file.txt")
println(key)
// Output: 2ojE41
```

#### Upload a file with an expiration date (here, 7 days)

```go
key, expiry, _ := fileio.UploadWithExpire("/data/file.txt", 7)
println(key + " expires in " + expiry)
// Output: aQbnDJ expires in 1 week
```

#### Download a file

```go
err := fileio.Download("2ojE41", "/tmp/file.txt")
```
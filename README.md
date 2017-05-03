# File.io

[![GoDoc](https://godoc.org/github.com/rvflash/fileio?status.svg)](https://godoc.org/github.com/rvflash/fileio)
[![Build Status](https://img.shields.io/travis/rvflash/fileio.svg)](https://travis-ci.org/rvflash/fileio)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/fileio.svg)](http://codecov.io/github/rvflash/fileio?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/fileio)](https://goreportcard.com/report/github.com/rvflash/fileio)


Golang interface for uploading or downloading files with file.io.


## Installation

```bash
$ go get -u github.com/rvflash/fileio
```

## Usage


#### Upload a file

```go
	key, err := fileio.Upload("/data/file.txt")
```

#### Upload a file with an expiration date (5 days in this example)

```go
	key, err := fileio.UploadWithExpire("/data/file.txt", 5)
```

#### Download a file

```go
	err := fileio.Download("2ojE41", "/tmp/file.txt")
```
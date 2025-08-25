# gurl

`gurl` is a lightweight Go implementation of the popular `curl` command-line tool. It allows you to make HTTP requests from the terminal using Go, supporting common features like GET, POST, headers, and more.

## Features

- Simple HTTP requests (GET, POST, PUT, DELETE)
- Custom headers and data
- Response output to terminal or file

## Installation

```sh
go install github.com/zer0go/gurl@latest
```

## Usage

```sh
gurl https://example.com
gurl -X POST -d '{"a":"b"}' -H 'content-type: application/json' https://postman-echo.com/post
```

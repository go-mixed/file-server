# File HTTP server

A simple HTTP server that serves files from the local file system.

It only supports GET requests, and HTML、CSS、JavaScript、image files. 

## Usage

```
$ ./file-server

# or
$ ./file-server -addr ":8080" -dir "/tmp"
```
- `-addr`: the address to listen on, default is ":8080"
- `-dir`: the directory to serve files from, default is "." (current directory)

## Compile

- Golang 1.18+
- no dependencies

```bash
$ cd file-server
$ go build
```

## Logging for stdout

```
127.0.0.1 - - [2023-10-19T12:00:42.1319939+08:00] "GET /main.go HTTP/1.1" 200 0B 139.6052ms "http://127.0.0.1:8080/main.go" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"
```
- client ip: the client's IP address
- username: the HTTP auth username, if any
- time: the RFC3339 formatted request time
- method: the HTTP method
- uri: the HTTP URI
- protocol: the HTTP protocol
- status: the HTTP status code
- size: the HTTP response size
- duration: the duration of processing the request
- referer: the HTTP referer
- user-agent: the HTTP user agent


## Alternative
- Python
```bash
python3 -m http.server 8080
python2 -m SimpleHTTPServer 8080
```

- PHP (>= 5.4.0)
```bash
php -S 0.0.0.0:8080
```

- Node.js (>= 8.0.0)
```bash
npm start -- --port=8080
```
或(>= 5.2.0)
```bash
npm install -g serve
serve -p 8080
```

- Ruby
```bash
ruby -run -e httpd . -p 8080
```
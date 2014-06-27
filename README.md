httpload-go
====

http_load like http(s) benchmark tool

## Usage

basic usage

```
$ echo http://localhost:5000/ > url_list.txt
$ httpload-go url_list.txt
option: parallel= 10 , seconds= 10
timeup!
result: score= 389.20 req/sec, statusCode=map[200:3892], total req count=3892
result: transfered: 81685296 bytes, mean: 8168529.60 bytes/sec, 20988.00 bytes/req
```

you can specify 2 options:

- `parallel` : loading worker counts
- `seconds` : loading time (in seconds)

also, url_list can multiple urls.

```
$ echo http://localhost:5000/foo >> url_list.txt
$ echo http://localhost:5000/bar >> url_list.txt
$ httpload-go --parallel 20 --seconds 5 url_list.txt
option: parallel= 10 , seconds= 5
timeup!
result: score= 391.60 req/sec, statusCode=map[200:1306 404:652], total req count=1958
result: transfered: 27086003 bytes, mean: 5417200.60 bytes/sec, 13833.51 bytes/req
```


## Author

[@acidlemon](https://twitter.com/acidlemon)


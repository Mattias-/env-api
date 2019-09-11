# env-api

A service exposing its environment variables.

## Usage

Run application

```bash
$ docker run --rm --network=host -e color="#ff0000" env-api:latest
2019/09/11 19:51:50 Starting env-api application...
127.0.0.1 - - [11/Sep/2019:19:51:53 +0000] "GET / HTTP/1.1" 200 99
```

Result

```bash
$ curl -s http://localhost:8080 | jq
{
  "HOME": "/",
  "HOSTNAME": "zzz",
  "PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
  "color": "#ff0000"
}
```

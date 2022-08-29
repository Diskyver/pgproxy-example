# pgproxy-example

A simple CLI postgresql server proxy made with this [pgproxy library](https://github.com/Diskyver/pgproxy).

```sh
go run *.go -h
usage:  main [options]
  -pg-addr string
        Postgresql url. Also the PGPROXY_DB_URL environment variable exists (default "postgres://127.0.0.1:5432")
  -proxy-addr string
        Proxy address (default "127.0.0.1:15432")
```

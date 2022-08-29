package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diskyver/pgproxy"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
)

var options struct {
	proxyAddr string
	pgAddr    string
}

func pgAddr() (string, error) {
	if addr := os.Getenv("PGPROXY_DB_URL"); addr != "" {
		return addr, nil
	}

	if addr := options.pgAddr; addr != "" {
		return addr, nil
	}

	return "", errors.New("postgresql url not provided")
}

type Session struct{}

func (s *Session) OnConnect(_ *pgx.Conn) {
	fmt.Println("I'm in!")
}

func (s *Session) OnQuery(query *pgproto3.Query) (*pgproto3.Query, error) {
	fmt.Println("query:", query.String)
	return query, nil
}

func (s *Session) OnResult(rows pgx.Rows, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("%w", rows)
}

func (s *Session) OnClose(_ *pgx.Conn) {
	fmt.Println("I'm out")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&options.proxyAddr, "proxy-addr", "127.0.0.1:15432", "Proxy address")
	flag.StringVar(&options.pgAddr, "pg-addr", "postgres://127.0.0.1:5432", "Postgresql url. Also the PGPROXY_DB_URL environment variable exists")
	flag.Parse()

	dbAddr, err := pgAddr()
	if err != nil {
		log.Fatal(err)
	}

	proxy := pgproxy.CreatePgProxy(dbAddr, &Session{})

	signal_channel := make(chan os.Signal)
	signal.Notify(signal_channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signal_channel
		fmt.Println("\rHandle SIGTERM")
		if err := proxy.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to gracefully shutdown the PgProxyServer: %e", err)
			os.Exit(1)
		}
		fmt.Println("PgProxyServer closed gracefully")
		os.Exit(0)
	}()

	proxy.Listen(options.proxyAddr)

}

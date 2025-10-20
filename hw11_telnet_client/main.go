package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "время ожадния соединения")
}

func main() {
	flag.Parse()

	address := flag.Arg(0)
	if address == "" {
		log.Fatal("Адрес подключения не указан")
		return
	}

	port := flag.Arg(1)
	if port == "" {
		log.Fatal("Порт подключения не указан")
		return
	}

	client := NewTelnetClient(net.JoinHostPort(address, port), timeout, os.Stdin, os.Stdout)

	if err := RunTelnetClient(client); err != nil {
		log.Fatal(err)
		return
	}
}

package main

import (
	"os"
	"os/signal"
	"syscall"
	"usermanagement/server"
)

//go:generate swagger generate spec

func main() {
	server := server.NewServer()
	go listenToSystemSignals(server)
	err := server.Run()
	code := server.ExitCode(err)
	os.Exit(code)
}

func listenToSystemSignals(server *server.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	for sig := range signalChan {
		reason := "system signal: " + sig.String()
		server.Shutdown(reason)
	}
}

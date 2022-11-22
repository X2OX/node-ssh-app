package main

import (
	"os"
	"os/signal"
	"syscall"

	nsa "github.com/X2OX/node-ssh-app"
	"go.uber.org/zap"
	"go.x2ox.com/THz"
	"go.x2ox.com/sorbifolia/rogu"
)

func init() {
	rogu.MustReplaceGlobals(rogu.DefaultZapConfig(rogu.DefaultZapEncoderConfig(),
		[]string{"stdout"},
		[]string{"stderr"}))
}

func main() {
	r := nsa.Router()
	r.SetLog(zap.L().Named("THz"))
	handleSignal(r)

	if err := r.ListenAndServe(":80"); err != nil {
		zap.L().Fatal("server exit", zap.Error(err))
	}
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *THz.THz) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		sig := <-c
		zap.L().Info("server exit signal", zap.Any("signal notify", sig))

		_ = server.Stop()
		os.Exit(0)
	}()
}

package helpers

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Starts goroutine for waiting signals and return channel with shutdown signal
func AwaitProcSignals() chan bool {
	// Init system signals
	shutdown := make(chan bool)
	//create a notification channel to shutdown
	sigChan := make(chan os.Signal, 1)
	//register for interupt (Ctrl+C) and SIGTERM (docker)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// Start goroutine for receiving signals
	go func() {
		// Wait any of signal in Signal Channel
		sig := <-sigChan
		logrus.Info("Signal received: ", sig.String())
		// Set shutdown = true
		shutdown <- true
	}()

	return shutdown
}

package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/kang-makes/cilium-nodeport-deleter/deleter"
)

func main() {
	// Flag and environment parsing.
	pflag.String("log-level", "info", "Log level (verbosity)")
	pflag.String("cilium-url", "unix:///var/run/cilium/cilium.sock", "URL where to talk to cilium socket")
	pflag.Parse()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("NODEPORT_DELETER")
	viper.AutomaticEnv()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err.Error())
	}

	// Logging
	lvl, _ := log.ParseLevel(viper.GetString("log-level"))
	log.SetLevel(lvl)

	// Signal management
	sigChan := make(chan os.Signal, 10)
	errChan := make(chan error)
	abortChan := make(chan interface{})
	defer close(sigChan)
	defer close(errChan)
	defer close(abortChan)

	go func() {
		nodePortDeleter, err := deleter.New(viper.GetString("cilium-url"))
		timer := time.NewTimer(10 * time.Second)
		if err != nil {
			errChan <- err
			return
		}

		for {
			err := nodePortDeleter.DeleteServices()
			if err != nil {
				log.Errorf("Error return by service deleter: %s", err)
			}

			select {
				case abort := <-abortChan:
					_ = abort
					return
				case <-timer.C:
			}

			timer.Reset(10 * time.Second)
		}
	}()

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case signal := <-sigChan:
		log.Printf("Caught %s, shutting down...", signal.String())

		abortChan <- nil

	case err := <-errChan:
		log.Errorf("Unrecoverable error while deleting services: %v", err)
	}
}

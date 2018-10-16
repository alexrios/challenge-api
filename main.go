package main

import (
	"github.com/alexrios/challenge-api/app"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	bootPhaseLogger := log.WithFields(log.Fields{
		"phase": "boot",
	})
	bootPhaseLogger.Info("Starting App...")
	//Criando App
	app, err := app.NewApp()
	if err != nil {
		bootPhaseLogger.WithFields(log.Fields{
			"cause": err.Error(),
		}).Fatal("Could not start app")
	}
	//Validando minimamente o endereco
	addr := os.Getenv("SERVER_ADDR")
	if len(addr) == 0 {
		bootPhaseLogger.WithFields(log.Fields{
			"cause": "Invalid address",
			"addr":  addr,
		}).Fatal("Could not start app.")
	}
	//Executando App
	app.Run(addr)
}

package main

import (
	"flag"
	"github.com/AleksK1NG/es-microservice/config"
	"github.com/AleksK1NG/es-microservice/internal/server"
	"github.com/AleksK1NG/es-microservice/pkg/logger"
	"log"
)

func main() {
	log.Println("Starting es microservice")

	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("EventSourcingService")

	appLogger.Infof("CFG: %+v", cfg)
	appLogger.Info("Success =D")
	appLogger.Fatal(server.NewServer(cfg, appLogger).Run())
}

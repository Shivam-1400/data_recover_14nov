package main

import (
	"context"
	"data_recover_14_nov/config"
	"data_recover_14_nov/databases"
	"data_recover_14_nov/globals"
	"data_recover_14_nov/services"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	fmt.Println("Hello, World!")
	configPath := "./config/config.json"

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	globals.ApplicationConfig = cfg
	fmt.Printf("Configuration loaded successfully: %+v\n", cfg)

	globals.DataMap = &sync.Map{}
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go services.FileRead(ctx, &wg)

	wg.Add(1)
	go services.CheckData(ctx, &wg)

	wg.Add(1)
	go databases.CheckConnection(ctx, &wg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	cancel()

	wg.Wait()
	fmt.Println("Application shutdown successfully.")
}

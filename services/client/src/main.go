package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	client "github.com/7574-sistemas-distribuidos/tp-nivelador/src/client"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
)

func loadConfig() (client.ClientConfig, error) {
	agencyId := os.Getenv("AGENCY_ID")
	if agencyId == "" {
		return client.ClientConfig{}, errors.New("AGENCY_ID environment variable is required")
	}

	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		return client.ClientConfig{}, errors.New("SERVER_HOST environment variable is required")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		return client.ClientConfig{}, errors.New("SERVER_PORT environment variable is required")
	}

	inputFilePath := os.Getenv("INPUT_FILE")
	if inputFilePath == "" {
		return client.ClientConfig{}, errors.New("INPUT_FILE environment variable is required")
	}

	outputFilePath := os.Getenv("OUTPUT_FILE")
	if outputFilePath == "" {
		return client.ClientConfig{}, errors.New("OUTPUT_FILE environment variable is required")
	}

	batchSize, err := loadBatchSize()
	if err != nil {
		return client.ClientConfig{}, err
	}

	return client.ClientConfig{
		ServerHost:     serverHost,
		ServerPort:     serverPort,
		AgencyID:       agencyId,
		InputFilePath:  inputFilePath,
		OutputFilePath: outputFilePath,
		BatchSize:      batchSize,
	}, nil
}

func loadBatchSize() (int, error) {
	batchSizeStr := os.Getenv("BATCH_SIZE")
	if batchSizeStr == "" {
		return client.DEFAULT_BATCH_SIZE, nil
	}

	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil || batchSize < 1 {
		return 0, fmt.Errorf("BATCH_SIZE must be a positive integer, got %q", batchSizeStr)
	}

	return batchSize, nil
}

func run() int {
	config, err := loadConfig()
	if err != nil {
		logger.Error("load-config", logger.Fail, "err", err)
		return 1
	}

	client, err := client.NewClient(config)
	if err != nil {
		logger.Error("client-new", logger.Fail, "err", err)
		return 1
	}

	if err := client.Run(); err != nil {
		logger.Error("client-run", logger.Fail, "err", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}

package client

import (
	"bufio"
	"net"
	"os"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

const (
	CONNECTION_ATTEMPTS_MAX     = 3
	CONNECTION_ATTEMPS_DELAY_MS = 200
)

const (
	ECHO_CLIENT_BUFFER_SIZE      = 512
	ECHO_CLIENT_MESSAGE_AMOUNT   = 3
	ECHO_CLIENT_MESSAGE_DELAY_MS = 1000
)

type ClientConfig struct {
	ServerHost     string
	ServerPort     string
	AgencyId       string
	InputFilePath  string
	OutputFilePath string
}

type Client struct {
	conn   net.Conn
	config ClientConfig
}

func NewClient(config ClientConfig) (*Client, error) {
	conn, err := connectToServer(config.ServerHost, config.ServerPort)
	if err != nil {
		logger.Warn("connect-to-server", logger.Fail)
		return nil, err
	}

	client := &Client{conn: conn, config: config}
	return client, nil
}

func connectToServer(host, port string) (net.Conn, error) {
	const action = "connect-to-server"
	var err error
	var conn net.Conn

	logger.Info(action, logger.InProgress)
	for i := range CONNECTION_ATTEMPTS_MAX {
		conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {
			logger.Warn(action, logger.Fail, "attempt", i)
			time.Sleep(CONNECTION_ATTEMPS_DELAY_MS * time.Millisecond)
			continue
		}

		logger.Info(action, logger.Success)
		break
	}

	return conn, err
}

func (client *Client) Run() error {
	defer client.conn.Close()
	const mainAction = "test-echo-server"
	messageId := 0

	inputFile, err := os.Open(client.config.InputFilePath)
	if err != nil {
		logger.Error("open-file", logger.Fail, err)
		return err
	}

	defer inputFile.Close()

	outputFile, err := os.Create(client.config.OutputFilePath)
	if err != nil {
		logger.Error("open-file", logger.Fail, err)
		return err
	}

	defer outputFile.Close()

	inputScanner := bufio.NewScanner(inputFile)
	outputWriter := bufio.NewWriter(outputFile)

	for inputScanner.Scan() {
		clientMessage := inputScanner.Text()
		messageArgs := []any{"agency-id", client.config.AgencyId, "message-id", messageId, "message", clientMessage}
		logger.Info(mainAction, logger.InProgress, messageArgs...)

		if err := safe_socket.SendAll(client.conn, []byte(clientMessage)); err != nil {
			logger.Error("send-message", logger.Fail, messageArgs...)
			return err
		}

		responseBuffer, err := safe_socket.RecvAll(client.conn, ECHO_CLIENT_BUFFER_SIZE)
		if err != nil {
			logger.Error("recv-response", logger.Fail, messageArgs...)
			return err
		}

		if string(responseBuffer) != clientMessage {
			messageArgs := []any{"agency-id", client.config.AgencyId, "message-id", messageId, "message", clientMessage, "buffer", string(responseBuffer)}
			logger.Error("check-response", logger.Fail, messageArgs...)
			return err
		}

		response := string(responseBuffer)
		_, err = outputWriter.WriteString(response + "\n")
		if err != nil {
			messageArgs := []any{"agency-id", client.config.AgencyId, "message-id", messageId, "message", clientMessage, "err", err}
			logger.Error("check-response", logger.Fail, messageArgs...)
			return err
		}
		outputWriter.Flush()
		messageId++
		time.Sleep(ECHO_CLIENT_MESSAGE_DELAY_MS * time.Millisecond)
	}

	logger.Info(mainAction, logger.Success, "agency-id", client.config.AgencyId)

	return nil
}

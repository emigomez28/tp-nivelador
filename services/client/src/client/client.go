package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol"
)

const (
	CONNECTION_ATTEMPTS_MAX     = 3
	CONNECTION_ATTEMPS_DELAY_MS = 200
	DEFAULT_BATCH_SIZE          = 8
)

type ClientConfig struct {
	ServerHost     string
	ServerPort     string
	AgencyID       string
	InputFilePath  string
	OutputFilePath string
	BatchSize      int
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
	action := "send-bets"

	logger.Info(action, logger.InProgress, "agency-id", client.config.AgencyID)

	if err := client.startTransmission(); err != nil {
		logger.Error(action, logger.Fail, "err", err)
		return err
	}

	betsAmount, err := client.sendBets()
	if err != nil {
		logger.Error(action, logger.Fail, "bets-amount", betsAmount, "err", err)
		return err
	}

	winners, err := client.endTransmission()
	if err != nil {
		logger.Error(action, logger.Fail, "err", err)
		return err
	}

	if err := client.storeWinners(winners); err != nil {
		logger.Error("store-winners", logger.Fail, "err", err)
		return err
	}

	logger.Info(action, logger.Success, "agency-id", client.config.AgencyID, "bets-amount", betsAmount)
	return nil
}

func (client *Client) startTransmission() error {
	msg := protocol.NewMessage(protocol.MsgStartTransmission, []byte(client.config.AgencyID))

	if err := protocol.SendMessage(client.conn, msg); err != nil {
		return err
	}

	return client.recvOk()
}

func (client *Client) sendBets() (int, error) {
	inputFile, err := os.Open(client.config.InputFilePath)
	if err != nil {
		return 0, err
	}
	defer inputFile.Close()

	betsAmount := 0
	batch := protocol.NewBetBatch(client.config.BatchSize)
	inputScanner := bufio.NewScanner(inputFile)

	for inputScanner.Scan() {
		betLine := inputScanner.Text()
		if !batch.CanAdd(betLine) {
			sent, err := client.sendBatch(batch)
			betsAmount += sent
			if err != nil {
				return betsAmount, err
			}
		}

		batch.Add(betLine)
	}

	if err := inputScanner.Err(); err != nil {
		return betsAmount, err
	}

	if batch.IsEmpty() {
		return betsAmount, nil
	}

	sent, err := client.sendBatch(batch)
	betsAmount += sent

	return betsAmount, err
}

func (client *Client) sendBatch(batch *protocol.BetBatch) (int, error) {
	betsAmount := batch.Count()
	err := protocol.SendMessage(client.conn, batch.Message())
	batch.Reset()

	if err != nil {
		return 0, err
	}

	if err := client.recvOk(); err != nil {
		return 0, err
	}

	return betsAmount, nil
}

func (client *Client) endTransmission() ([]byte, error) {
	msg := protocol.NewMessage(protocol.MsgEndTransmission, []byte{})
	if err := protocol.SendMessage(client.conn, msg); err != nil {
		return nil, err
	}

	response, err := client.recvExpecting(protocol.MsgWinners)
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func (client *Client) recvOk() error {
	_, err := client.recvExpecting(protocol.MsgOk)
	return err
}

func (client *Client) recvExpecting(expected protocol.MessageType) (*protocol.Message, error) {
	response, err := protocol.RecvMessage(client.conn)
	if err != nil {
		return nil, err
	}

	if response.Type != expected {
		return nil, fmt.Errorf("expected %s, got %s", expected, response.Type)
	}

	return response, nil
}

func (client *Client) storeWinners(winners []byte) error {
	outputFile, err := os.Create(client.config.OutputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if len(winners) == 0 {
		return nil
	}

	if _, err := outputFile.Write(winners); err != nil {
		return err
	}

	_, err = outputFile.WriteString("\n")
	return err
}

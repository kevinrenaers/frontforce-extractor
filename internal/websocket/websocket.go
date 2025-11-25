package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WebsocketClient struct {
	baseUrl      string
	messagesChan chan json.RawMessage

	conn *websocket.Conn
}

type message struct {
	Type      int             `json:"type"`
	Target    string          `json:"target"`
	Arguments json.RawMessage `json:"arguments"`
}

func NewWebSocketClient(url string) WebsocketClient {
	return WebsocketClient{
		baseUrl:      url,
		messagesChan: make(chan json.RawMessage),
	}
}

func (w WebsocketClient) Connect(connectionToken, accessToken string) error {
	url := fmt.Sprintf("%s&id=%s&access_token=%s", w.baseUrl, connectionToken, accessToken)

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Printf("Dial failed: %s\n", err.Error())

		return fmt.Errorf("failed setting up websocket client: %w", err)
	}

	rawMessage := `{"protocol":"json","version":1}`
	err = ws.WriteMessage(websocket.TextMessage, []byte(rawMessage))
	if err != nil {
		return fmt.Errorf("failed sending initial websocket message: %w", err)
	}

	w.conn = ws

	go func() {
		for {
			if w.conn == nil {
				return
			}

			var message []byte

			_, message, err := w.conn.ReadMessage()
			if err != nil {
				log.Error().Err(err).Msg("websocket - error reading websocket message")

				return
			}

			w.parseMessage(message)
		}
	}()

	log.Info().Msg("websocket - successfully started websocket connection")

	return nil
}

func (w WebsocketClient) Close() error {
	if w.conn == nil {
		return nil
	}
	err := w.conn.Close()
	if err != nil {
		return fmt.Errorf("failed closing websocket connection: %w", err)
	}

	w.conn = nil

	log.Info().Msg("websocket - successfully closed websocket connection")

	return nil
}

func (w WebsocketClient) Messages() chan json.RawMessage {
	return w.messagesChan
}

func (w WebsocketClient) parseMessage(data []byte) {
	if len(data) <= 3 {
		return
	}

	data = data[:len(data)-1] // Remove the trailing '' character

	var msg message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Error().Err(err).Msg("websocket - failed parsing message")

		return
	}

	switch msg.Type {
	case 1:
		log.Info().Msgf("websocket - received message type: %d", msg.Type)
		w.messagesChan <- msg.Arguments
	default:
		log.Info().Msgf("websocket - received message type: %d", msg.Type)
	}
}

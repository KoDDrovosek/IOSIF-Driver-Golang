package IOSIF_Driver_Golang

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Connector interface {
	Pull(topicId string) (key string, value json.RawMessage, err error)
	Subscribe(topic string, handler func(key string, value json.RawMessage)) error
	BulkSubscribe(topics map[string]func(key string, value json.RawMessage))
	Publish(topic, key, value string) error
	Start() error
}

type message struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type Options struct {
	Topics  map[string]func(key string, value json.RawMessage)
	URL     string
	Periods int64
}

type SubscribtionResponse struct {
	Token string `json:"token"`
}

type connector struct {
	subscriberId string
	isListenerUp bool
	Options
}

func (c *connector) BulkSubscribe(topics map[string]func(key string, value json.RawMessage)) {
	c.Topics = topics
}

func (c *connector) Subscribe(topic string, handler func(key string, value json.RawMessage)) error {
	response, err := http.Post(fmt.Sprintf("%s/subscribe", c.URL), "application/json", bytes.NewBufferString(fmt.Sprintf("[\"%s\"]", topic)))
	if err != nil || response.StatusCode != http.StatusCreated {
		return err
	}
	var marshaledResponse SubscribtionResponse
	json.NewDecoder(response.Body).Decode(&marshaledResponse)
	c.subscriberId = marshaledResponse.Token
	c.Topics[topic] = handler
	return nil
}

func (c connector) Publish(topicId, key, value string) error {

	payload := fmt.Sprintf("{\"key\":\"%s\", \"value\": %s}", key, value)
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/publish", c.URL), bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}

	q := request.URL.Query()
	q.Set("topicId", topicId)
	request.URL.RawQuery = q.Encode()

	client := http.Client{}

	response, err := client.Do(request)

	if response.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("server response with status code %d", response.StatusCode))
	}

	return nil
}

func (c connector) Pull(topicId string) (key string, value json.RawMessage, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/pull", c.URL), nil)
	if err != nil {
		return "", nil, err
	}

	q := req.URL.Query()
	q.Set("topicId", topicId)
	q.Set("subscriberId", c.subscriberId)
	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return "", nil, errors.New("No messages")
	}

	var m message
	err = json.NewDecoder(res.Body).Decode(&m)
	return m.Key, m.Value, err
}

func (c connector) Start() error {
	if c.isListenerUp {
		return errors.New("listener is already up")
	}

	go c.listener()

	return nil
}

func (c *connector) listener() {
	for {
		for topic, handler := range c.Topics {
			key, value, err := c.Pull(topic)
			if err != nil {
				continue
			}

			handler(key, value)
		}

		time.Sleep(time.Duration(c.Periods))
	}
}

func New(opt Options) Connector {
	return &connector{
		subscriberId: "",
		isListenerUp: false,
		Options:      opt,
	}
}

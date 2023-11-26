package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type Params struct {
	DefaultInstance string
	DefaultTopic    string
	Token           string
}

type Publisher interface {
	Publish(ctx context.Context, publication *Publication) error
}

type Publication struct {
	Title              string
	Topic              string
	Priority           int
	Message            string
	MessageIsMarkdown  bool `json:"markdown"`
	Tags               []string
	ClickURL           string `json:"click"`
	Delay              string
	Email              string
	Call               string
	AttachmentURL      string `json:"attach"`
	AttachmentFilename string `json:"filename"`
	IconURL            string `json:"icon"`
	InstanceURL        string `json:"-"`
	Token              string `json:"-"`
	Actions            []PublicationAction
}

type PublicationAction struct {
	Action string
	Label  string
	Clear  bool
	// for "view" and "http" action only
	URL string
	// for "broadcast" action only
	Extras map[string]string
	// for "broadcast" action only
	Intent string
	// for "http" action only
	Method string
	// for "http" action only
	Headers map[string]string
	// for "http" action only
	Body string
}

func (p *Publication) MarshalZerologObject(e *zerolog.Event) {
	e.Int("priority", p.Priority).
		Str("title", p.Title).
		Str("message", p.Message)
}

type NtfyPublisher struct {
	client *http.Client
	params Params
}

func NewNtfyPublisher(params Params) *NtfyPublisher {
	return &NtfyPublisher{
		params: params,
		client: &http.Client{},
	}
}

func (p *NtfyPublisher) Publish(ctx context.Context, publication *Publication) error {
	payload, err := json.Marshal(publication)
	if err != nil {
		return err
	}

	instance, topic := p.getInstanceAndTopic(publication)
	req, err := http.NewRequest("POST", instance, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	token := publication.Token
	if token == "" {
		token = p.params.Token
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to publish to topic '%s' on '%s': %s", topic, instance, body)
	}

	return err
}

func (p *NtfyPublisher) getInstanceAndTopic(publication *Publication) (instance, topic string) {
	topic = publication.Topic
	if topic == "" {
		topic = p.params.DefaultTopic
	}
	instance = publication.InstanceURL
	if instance == "" {
		instance = p.params.DefaultInstance
	}
	return
}

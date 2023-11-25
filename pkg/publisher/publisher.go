package publisher

import (
	"heckel.io/ntfy/v2/client"
)

type Params struct {
	Token string
}

type Publisher interface {
	Publish(topic string, publication *Publication) error
}

type Publication struct {
	Priority string
	Title    string
	Message  string
}

func (p *Publication) PublishOptions() []client.PublishOption {
	opts := []client.PublishOption{}
	if p.Priority != "" {
		opts = append(opts, client.WithPriority(p.Priority))
	}
	if p.Title != "" {
		opts = append(opts, client.WithTitle(p.Title))
	}
	if p.Message != "" {
		opts = append(opts, client.WithMessage(p.Message), client.WithMarkdown())
	}
	return opts
}

type NtfyPublisher struct {
	client *client.Client
	params Params
}

func NewNtfyPublisher(params Params) *NtfyPublisher {
	return &NtfyPublisher{
		params: params,
		client: client.New(&client.Config{}),
	}
}

func (p *NtfyPublisher) Publish(topic string, publication *Publication) error {
	opts := append(publication.PublishOptions(), client.WithBearerAuth(p.params.Token))

	p.client.Publish(topic, "", opts...)
	return nil
}

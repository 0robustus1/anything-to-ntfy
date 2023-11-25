package slack

import (
	"fmt"

	"github.com/0robustus1/anything-to-ntfy/pkg/publisher"
	"github.com/gofiber/fiber/v2"
)

type Params struct {
	Publisher    publisher.Publisher
	DefaultTopic string
}

type SlackInput struct {
	params Params
}

func NewSlackInput(params Params) *SlackInput {
	return &SlackInput{
		params: params,
	}
}

func (i *SlackInput) RegisterWithRouter(router fiber.Router) {
	router.Post("/slack/incoming_webhook/:topic", i.handleIncomingWebhook)
}

type slackIncomingWebhookMessage struct {
	Text string `json:"text"`
}

func (m *slackIncomingWebhookMessage) Publication() *publisher.Publication {
	return &publisher.Publication{
		Message: m.Text,
	}
}

func (i *SlackInput) handleIncomingWebhook(c *fiber.Ctx) error {
	topic := c.Params("topic", i.params.DefaultTopic)
	if topic == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Topic parameter must be provided")
	}

	message := &slackIncomingWebhookMessage{}
	if err := c.BodyParser(message); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to process payload from slack incoming webhook: %v", err))
	}

	if err := i.params.Publisher.Publish(topic, message.Publication()); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to publish message from slack incoming webhook: %v", err))
	}

	return nil
}

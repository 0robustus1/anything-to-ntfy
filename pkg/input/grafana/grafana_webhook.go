package grafana

import (
	"fmt"

	"github.com/0robustus1/anything-to-ntfy/pkg/publisher"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Params struct {
	Publisher publisher.Publisher
}

type GrafanaInput struct {
	params Params
}

func NewGrafanaInput(params Params) *GrafanaInput {
	return &GrafanaInput{
		params: params,
	}
}

func (i *GrafanaInput) RegisterWithRouter(router fiber.Router) {
	group := router.Group("/grafana")
	group.Post("/webhook", i.handleWebhook)
	group.Post("/webhook/:topic", i.handleWebhook)
}

type grafanaWebhookMessage struct {
	Status      string
	Alerts      []grafanaWebhookAlert `json:"alerts"`
	ExternalURL string                `json:"externalURL"`
}

type grafanaWebhookAlert struct {
	Status       string
	Labels       map[string]string
	Annotations  map[string]string
	StartsAt     string
	EndsAt       string
	GeneratorURL string
	SilenceURL   string
	DashboardURL string
	PanelURL     string
	Fingerprint  string
	Values       map[string]interface{}
}

func (i *GrafanaInput) handleWebhook(c *fiber.Ctx) error {
	message := &grafanaWebhookMessage{}
	if err := c.BodyParser(message); err != nil {
		log.Ctx(c.UserContext()).Err(err).Msg("failed to process payload from grafana webhook")
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to process payload from grafana webhook: %v", err))
	}

	for _, alert := range message.Alerts {
		title := fmt.Sprintf("[%s] %s: %s", alert.Status, alert.Labels["alertname"], alert.Annotations["summary"])
		pub := &publisher.Publication{
			Topic:       c.Params("topic"),
			InstanceURL: c.Query("ntfyInstance"),
			Token:       c.Query("ntfyToken"),
			Title:       title,
			Message:     alert.Annotations["description"],
			Actions:     []publisher.PublicationAction{},
		}

		if alert.GeneratorURL != "" {
			pub.Actions = append(pub.Actions, publisher.PublicationAction{
				Action: "view",
				Label:  "Edit Alert",
				URL:    alert.GeneratorURL,
			})
		}

		if alert.SilenceURL != "" {
			pub.Actions = append(pub.Actions, publisher.PublicationAction{
				Action: "view",
				Label:  "Silence Alert",
				URL:    alert.SilenceURL,
			})
		}

		if alert.DashboardURL != "" {
			pub.Actions = append(pub.Actions, publisher.PublicationAction{
				Action: "view",
				Label:  "Dashboard",
				URL:    alert.DashboardURL,
			})
		}

		if alert.PanelURL != "" {
			pub.Actions = append(pub.Actions, publisher.PublicationAction{
				Action: "view",
				Label:  "Panel",
				URL:    alert.PanelURL,
			})
		}

		if err := i.params.Publisher.Publish(c.UserContext(), pub); err != nil {
			log.Ctx(c.UserContext()).Err(err).Str("topic", pub.Topic).Object("publication", pub).Msg("failed to publish message")
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to publish message from grafana webhook: %v", err))
		}

		log.Ctx(c.UserContext()).Info().Str("topic", pub.Topic).Object("publication", pub).Msg("published message")
	}

	return nil
}

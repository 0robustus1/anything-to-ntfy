package alertmanager

import (
	"fmt"
	"strings"

	"github.com/0robustus1/anything-to-ntfy/pkg/input"
	"github.com/0robustus1/anything-to-ntfy/pkg/publisher"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Params struct {
	Publisher publisher.Publisher
}

type AlertmanagerInput struct {
	params Params
}

func NewAlertmanagerInput(params Params) *AlertmanagerInput {
	return &AlertmanagerInput{
		params: params,
	}
}

func (i *AlertmanagerInput) RegisterWithRouter(router fiber.Router) {
	group := router.Group("/alertmanager")
	group.Post("/webhook", i.handleWebhook)
	group.Post("/webhook/:topic", i.handleWebhook)
}

func (i *AlertmanagerInput) handleWebhook(c *fiber.Ctx) error {
	ntfyInfo := input.NtfyInfoFromFiberContext(c)
	if err := ntfyInfo.Validate(); err != nil {
		log.Ctx(c.UserContext()).Err(err).Str("topic", ntfyInfo.Topic).Msg("invalid explicit ntfy config provided")
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to publish message from alertmanager webhook: %v", err))
	}
	message := &webhook.Message{}
	if err := c.BodyParser(message); err != nil {
		log.Ctx(c.UserContext()).Err(err).Msg("failed to process payload from alertmanager webhook")
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to process payload from alertmanager webhook: %v", err))
	}

	publications := publicationsFromMessage(message)
	for _, publication := range publications {
		if err := i.params.Publisher.Publish(c.UserContext(), publication, ntfyInfo); err != nil {
			log.Ctx(c.UserContext()).Err(err).Str("topic", ntfyInfo.Topic).Object("publication", publication).Msg("failed to publish message")
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to publish message from alertmanager webhook: %v", err))
		}

		log.Ctx(c.UserContext()).Info().Str("topic", ntfyInfo.Topic).Object("publication", publication).Msg("published message")
	}

	return nil
}

func statusTag(status string) (tag string) {
	if status == "resolved" {
		return "white_check_mark"
	}
	return "rotating_light"
}

func publicationsFromMessage(message *webhook.Message) []*publisher.Publication {
	publications := []*publisher.Publication{}
	c := cases.Title(language.AmericanEnglish)

	for _, alert := range message.Alerts {
		title := fmt.Sprintf("[%s] %s: %s", alert.Status, alert.Labels["alertname"], alert.Annotations["summary"])
		pub := &publisher.Publication{
			Title:   title,
			Tags:    []string{statusTag(alert.Status), alert.Labels["severity"]},
			Message: alert.Annotations["description"],
			// Note only up to 3 actions are allowed
			Actions: []publisher.PublicationAction{},
		}

		labelKey, url := playbookURL(&alert)
		if labelKey != "" {
			prettyLabel := c.String(strings.TrimSuffix(labelKey, "_url"))
			pub.Actions = append(pub.Actions, publisher.PublicationAction{
				Action: "view",
				Label:  prettyLabel,
				URL:    url,
			})
		}

		publications = append(publications, pub)
	}

	return publications
}

var playbookKeys = []string{"playbook_url", "runbook_url"}

func playbookURL(alert *template.Alert) (key, url string) {
	for _, checkKey := range playbookKeys {
		if val, ok := alert.Labels[checkKey]; ok {
			key = checkKey
			url = val
			return
		}
	}
	return
}

package input

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type NtfyInfo struct {
	Topic       string
	InstanceURL string
	Token       string
}

func NtfyInfoFromFiberContext(c *fiber.Ctx) *NtfyInfo {
	topic := c.Params("topic")
	if topic == "" {
		topic = c.Query("ntfyTopic")
	}
	if topic == "" {
		topic = c.Get("X-NTFY-Topic")
	}

	instanceURL := c.Query("ntfyInstance")
	if instanceURL == "" {
		instanceURL = c.Get("X-NTFY-INSTANCE")
	}

	token := c.Query("ntfyToken")
	if token == "" {
		token = c.Get("X-NTFY-TOKEN")
	}

	return &NtfyInfo{
		Topic:       topic,
		InstanceURL: instanceURL,
		Token:       token,
	}
}

func (i *NtfyInfo) Validate() error {
	if i.InstanceURL != "" && i.Token == "" {
		return fmt.Errorf("customizing ntfy instance requires setting a token as well")
	}
	return nil
}

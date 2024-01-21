package main

import (
	"fmt"

	"github.com/0robustus1/anything-to-ntfy/pkg/input/alertmanager"
	"github.com/0robustus1/anything-to-ntfy/pkg/input/grafana"
	"github.com/0robustus1/anything-to-ntfy/pkg/input/slack"
	"github.com/0robustus1/anything-to-ntfy/pkg/publisher"
	"github.com/alecthomas/kong"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var CLI struct {
	NtfyToken           string `env:"NTFY_TOKEN" required:"true" help:"Token to use to communicate with ntfy instance"`
	NtfyDefaultInstance string `env:"NTFY_DEFAULT_INSTANCE" optional:"" help:"Which ntfy instance to use by default"`
	NtfyDefaultTopic    string `env:"NTFY_DEFAULT_TOPIC" optional:"" help:"Which ntfy topic to use by default"`
	ListenHost          string `env:"LISTEN_HOST" optional:"" help:"Which host to listen on, should be an address. Defaults to empty string which is equivalent to 0.0.0.0"`
	ListenPort          int    `env:"LISTEN_PORT" optional:"" default:"5000" help:"Which port to listen on."`
}

func main() {
	_ = kong.Parse(&CLI)
	publisher := publisher.NewNtfyPublisher(publisher.Params{
		DefaultInstance: CLI.NtfyDefaultInstance,
		DefaultTopic:    CLI.NtfyDefaultTopic,
		Token:           CLI.NtfyToken,
	})
	slackInput := slack.NewSlackInput(slack.Params{
		Publisher: publisher,
	})
	grafanaInput := grafana.NewGrafanaInput(grafana.Params{
		Publisher: publisher,
	})
	alertmanagerInput := alertmanager.NewAlertmanagerInput(alertmanager.Params{
		Publisher: publisher,
	})

	app := fiber.New()
	logger := log.Logger.With().Str("app", "anything-to-ntfy").Logger()
	app.Use(func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		ctx = logger.WithContext(ctx)
		c.SetUserContext(ctx)
		return c.Next()
	})
	slackInput.RegisterWithRouter(app)
	grafanaInput.RegisterWithRouter(app)
	alertmanagerInput.RegisterWithRouter(app)
	app.Listen(fmt.Sprintf("%s:%d", CLI.ListenHost, CLI.ListenPort))
}

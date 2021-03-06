package pkg

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/krok-o/terraform-provider-krok/pkg/clients"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/auth"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/command"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/event"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/platform"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/repository"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/runs"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/setting"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/user"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/vault"
	"github.com/krok-o/terraform-provider-krok/pkg/clients/vcs"
)

// Config defines configuration for the Krok server client.
type Config struct {
	APIKeyID     string
	APIKeySecret string
	Address      string
	Email        string
}

// KrokClient is the main client for the Krok server.
// TODO: Break these down into interfaces and multiple clients.
type KrokClient struct {
	ApiKeyClient     *auth.Client
	CommandClient    *command.Client
	CommandRunClient *runs.Client
	EventClient      *event.Client
	PlatformClient   *platform.Client
	RepositoryClient *repository.Client
	SettingsClient   *setting.Client
	UserClient       *user.Client
	VaultClient      *vault.Client
	VcsClient        *vcs.Client
}

// NewKrokClient creates a new Krok server client.
func NewKrokClient(cfg Config, log zerolog.Logger) *KrokClient {
	handler := clients.NewHandler(clients.Config{
		APIKeyID:     cfg.APIKeyID,
		APIKeySecret: cfg.APIKeySecret,
		Address:      cfg.Address,
		Email:        cfg.Email,
		Client:       http.DefaultClient,
		Logger:       log,
	})
	apiKeyClient := auth.NewClient(cfg.Address, log, handler)
	commandClient := command.NewClient(cfg.Address, log, handler)
	commandRunClient := runs.NewClient(cfg.Address, log, handler)
	eventsClient := event.NewClient(cfg.Address, log, handler)
	platformClient := platform.NewClient(cfg.Address, log, handler)
	repoClient := repository.NewClient(cfg.Address, log, handler)
	settingsClient := setting.NewClient(cfg.Address, log, handler)
	userClient := user.NewClient(cfg.Address, log, handler)
	vaultClient := vault.NewClient(cfg.Address, log, handler)
	vcsClient := vcs.NewClient(cfg.Address, log, handler)
	return &KrokClient{
		ApiKeyClient:     apiKeyClient,
		CommandClient:    commandClient,
		CommandRunClient: commandRunClient,
		EventClient:      eventsClient,
		PlatformClient:   platformClient,
		RepositoryClient: repoClient,
		SettingsClient:   settingsClient,
		UserClient:       userClient,
		VaultClient:      vaultClient,
		VcsClient:        vcsClient,
	}
}

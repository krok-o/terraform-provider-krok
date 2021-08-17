package krok

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rs/zerolog"

	"github.com/krok-o/terraform-provider-krok/pkg"
)

// Provider defines an Krok Terraform provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KROK_API_KEY_ID", ""),
				Description: "KROK API KEY ID",
			},
			"api_key_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KROK_API_KEY_SECRET", ""),
				Description: "KROK API KEY SECRET",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KROK_EMAIL", ""),
				Description: "KROK EMAIL",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KROK_ENDPOINT", "http://localhost:9998"),
				Description: "KROK API ENDPOINT",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"krok_repository": resourceRepository(),
			"krok_command":    resourceCommand(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"krok_command":   dataSourceKrokCommand(),
			"krok_platform":  dataSourceKrokPlatform(),
			"krok_platforms": dataSourceKrokPlatforms(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Set up the main client.
	log := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}).With().Timestamp().Logger()
	client := pkg.NewKrokClient(pkg.Config{
		Address:      d.Get("endpoint").(string),
		APIKeyID:     d.Get("api_key_id").(string),
		APIKeySecret: d.Get("api_key_secret").(string),
		Email:        d.Get("email").(string),
	}, log)

	return client, nil
}

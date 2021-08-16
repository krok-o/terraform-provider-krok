package krok

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/krok-o/krok/pkg/models"

	"github.com/krok-o/terraform-provider-krok/pkg"
)

const (
	// platform data source fields
	platformIdFieldName   = "id"
	platformNameFieldName = "name"
)

// dataSourceKrokPlatform defines a Platform datasource terraform type.
func dataSourceKrokPlatform() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKrokPlatformRead,

		Schema: map[string]*schema.Schema{
			platformIdFieldName: {
				Type:     schema.TypeInt,
				Required: true,
			},
			platformNameFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// dataSourceKrokPlatformRead reloads the resource object from the terraform store.
func dataSourceKrokPlatformRead(data *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	cid := data.Get(platformIdFieldName).(int)
	platform, err := client.PlatformClient.Get(cid)
	if err != nil {
		return err
	}

	for k, v := range flattenPlatformObject(platform) {
		if err := data.Set(k, v); err != nil {
			return err
		}
	}
	data.SetId(strconv.Itoa(platform.ID))
	return nil
}

// flattenPlatformObject creates a map from an Krok Platform for easy digestion by the terraform schema.
func flattenPlatformObject(platform *models.Platform) map[string]interface{} {
	return map[string]interface{}{
		platformIdFieldName:   platform.ID,
		platformNameFieldName: platform.Name,
	}
}

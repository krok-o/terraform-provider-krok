package krok

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/krok-o/krok/pkg/models"

	"github.com/krok-o/terraform-provider-krok/pkg"
)

const (
	// platform data source fields
	platformsPlatformsFieldName     = "platforms"
	platformsPlatformsIdFieldName   = "id"
	platformsPlatformsNameFieldName = "name"
)

// dataSourceKrokPlatform defines a Platform datasource terraform type.
func dataSourceKrokPlatforms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKrokPlatformsRead,

		Schema: map[string]*schema.Schema{
			platformsPlatformsFieldName: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						platformsPlatformsIdFieldName: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						platformsPlatformsNameFieldName: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceKrokPlatformsRead reloads the resource object from the terraform store.
func dataSourceKrokPlatformsRead(data *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	platforms, err := client.PlatformClient.List()
	if err != nil {
		return err
	}

	for k, v := range flattenPlatformsObject(platforms) {
		if err := data.Set(k, v); err != nil {
			return err
		}
	}
	data.SetId(uniqueResourceID())
	return nil
}

// flattenPlatformsObject creates a map from an Krok Platform for easy digestion by the terraform schema.
func flattenPlatformsObject(platforms []models.Platform) map[string]interface{} {
	ret := make([]interface{}, 0)
	for _, v := range platforms {
		ret = append(ret, map[string]interface{}{
			platformsPlatformsIdFieldName:   v.ID,
			platformsPlatformsNameFieldName: v.Name,
		})
	}
	return map[string]interface{}{
		platformsPlatformsFieldName: ret,
	}
}

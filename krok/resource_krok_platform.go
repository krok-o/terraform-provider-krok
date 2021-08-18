package krok

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/krok-o/krok/pkg/models"

	"github.com/krok-o/terraform-provider-krok/pkg"
)

const (
	platformTokenFieldName = "token"
	platformVCSFieldName   = "vcs"
)

func resourcePlatform() *schema.Resource {
	return &schema.Resource{
		Create: resourcePlatformCreate,
		Read:   resourcePlatformRead,
		Update: resourcePlatformUpdate,
		Delete: resourcePlatformDelete,

		Schema: map[string]*schema.Schema{
			platformTokenFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			platformVCSFieldName: {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

// resourcePlatformCreate creates a Krok platform.
func resourcePlatformCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	expandedVCSToken, err := expandVCSTokenResource(d)
	if err != nil {
		return err
	}
	if err := client.VcsClient.Create(expandedVCSToken); err != nil {
		log.Println("Failed to create vcstoken.")
		return fmt.Errorf("failed to create vcstoken: %w", err)
	}
	d.SetId(uniqueResourceID())
	return resourceCommandRead(d, m)
}

// expandVCSTokenResource creates a Krok vcstoken structure out of a Terraform schema model.
func expandVCSTokenResource(d *schema.ResourceData) (*models.VCSToken, error) {
	var (
		token string
		vcs   int
	)
	if v, ok := d.GetOk(platformTokenFieldName); ok {
		token = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", platformTokenFieldName)
	}
	if v, ok := d.GetOk(platformVCSFieldName); ok {
		vcs = v.(int)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", platformVCSFieldName)
	}
	platform := &models.VCSToken{
		Token: token,
		VCS:   vcs,
	}
	return platform, nil
}

// resourcePlatformUpdate updates platform information from terraform stores.
func resourcePlatformUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

// resourcePlatformRead retrieves platform information from terraform stores.
func resourcePlatformRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePlatformDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

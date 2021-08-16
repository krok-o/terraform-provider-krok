package krok

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/krok-o/krok/pkg/models"

	"github.com/krok-o/terraform-provider-krok/pkg"
)

const (
	// command data source fields
	commandIdFieldName   = "id"
	commandNameFieldName = "name"
)

// dataSourceKrokCommand defines a Command datasource terraform type.
func dataSourceKrokCommand() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKrokCommandRead,

		Schema: map[string]*schema.Schema{
			commandIdFieldName: {
				Type:     schema.TypeInt,
				Required: true,
			},
			commandNameFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// dataSourceKrokCommandRead reloads the resource object from the terraform store.
func dataSourceKrokCommandRead(data *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	cid := data.Get(commandIdFieldName).(int)
	command, err := client.CommandClient.Get(cid)
	if err != nil {
		return err
	}

	for k, v := range flattenCommandObject(command) {
		if err := data.Set(k, v); err != nil {
			return err
		}
	}
	data.SetId(strconv.Itoa(command.ID))
	return nil
}

// flattenCommandObject creates a map from an Krok Command for easy digestion by the terraform schema.
func flattenCommandObject(command *models.Command) map[string]interface{} {
	return map[string]interface{}{
		commandIdFieldName:   command.ID,
		commandNameFieldName: command.Name,
	}
}

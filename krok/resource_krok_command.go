package krok

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/krok-o/krok/pkg/models"

	"github.com/krok-o/terraform-provider-krok/pkg"
)

const (
	commandResourceNameFieldName         = "name"
	commandResourceImageFieldName        = "image"
	commandResourceScheduleFieldName     = "schedule"
	commandResourcePlatformsFieldName    = "platforms"
	commandResourceRepositoriesFieldName = "repositories"
	commandResourceEnabledFieldName      = "enabled"
	// repoEventsFieldName                      = "name"
)

func resourceCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceCommandCreate,
		Read:   resourceCommandRead,
		Update: resourceCommandUpdate,
		Delete: resourceCommandDelete,

		Schema: map[string]*schema.Schema{
			commandResourceNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			commandResourceImageFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			commandResourceScheduleFieldName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			commandResourceEnabledFieldName: {
				Type:     schema.TypeBool,
				Required: true,
			},
			// TODO: Add support for updating this field.
			commandResourcePlatformsFieldName: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			commandResourceRepositoriesFieldName: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

// resourceCommandCreate creates a Krok repository.
func resourceCommandCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	expandedCommand, err := expandCommandResource(d)
	if err != nil {
		return err
	}
	createdCommand, err := client.CommandClient.Create(expandedCommand)
	if err != nil {
		log.Println("Failed to create command.")
		return fmt.Errorf("failed to create command: %w", err)
	}

	// add any relationships that might exist for commands.
	if v, ok := d.GetOk(commandResourcePlatformsFieldName); ok {
		providers := v.([]interface{})
		for _, pid := range providers {
			if err := client.CommandClient.AddRelationshipToPlatform(createdCommand.ID, pid.(int)); err != nil {
				log.Println("Failed to create relationship for command and platform.")
				return fmt.Errorf("failed to add relationship between command %d and platform %d: %w", createdCommand.ID, pid, err)
			}
		}
	}
	d.SetId(strconv.Itoa(createdCommand.ID))
	return resourceCommandRead(d, m)
}

// expandCommandResource creates a Krok command structure out of a Terraform schema model.
func expandCommandResource(d *schema.ResourceData) (*models.Command, error) {
	var (
		name     string
		image    string
		schedule string
		enabled  bool
	)
	if v, ok := d.GetOk(commandResourceNameFieldName); ok {
		name = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find or parse field %s", commandResourceNameFieldName)
	}
	if v, ok := d.GetOk(commandResourceImageFieldName); ok {
		image = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find or parse field %s", commandResourceImageFieldName)
	}
	if v, ok := d.GetOk(commandResourceEnabledFieldName); ok {
		enabled = v.(bool)
	} else {
		return nil, fmt.Errorf("unable to find or parse field %s", commandResourceEnabledFieldName)
	}
	if v, ok := d.GetOk(commandResourceScheduleFieldName); ok {
		schedule = v.(string)
	}
	command := &models.Command{
		Name:     name,
		Image:    image,
		Schedule: schedule,
		Enabled:  enabled,
	}
	return command, nil
}

// resourceCommandRead retrieves command information from terraform stores.
func resourceCommandRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	cid, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	command, err := client.CommandClient.Get(cid)
	if err != nil {
		log.Println("Failed to find command")
		d.SetId("")
		return err
	}

	for k, v := range flattenCommand(command) {
		if err := d.Set(k, v); err != nil {
			d.SetId("")
			return err
		}
	}
	return nil
}

// flattenCommand creates a map from a command for easy storage on terraform.
func flattenCommand(command *models.Command) map[string]interface{} {
	var (
		repositories []int
		platforms    []int
	)
	for _, r := range command.Repositories {
		repositories = append(repositories, r.ID)
	}
	for _, r := range command.Platforms {
		platforms = append(platforms, r.ID)
	}
	flatCommand := map[string]interface{}{
		commandResourceNameFieldName:         command.Name,
		commandResourceImageFieldName:        command.Image,
		commandResourceScheduleFieldName:     command.Schedule,
		commandResourceEnabledFieldName:      command.Enabled,
		commandResourcePlatformsFieldName:    platforms,
		commandResourceRepositoriesFieldName: repositories,
	}
	return flatCommand
}

// resourceCommandUpdate checks fields for differences and updates a command if necessary.
func resourceCommandUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	cid, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	command, err := client.CommandClient.Get(cid)
	if err != nil {
		log.Println("Failed to find command")
		d.SetId("")
		return err
	}

	if d.HasChange(commandResourceNameFieldName) {
		command.Name = d.Get(commandResourceNameFieldName).(string)
	}

	if d.HasChange(commandResourceScheduleFieldName) {
		command.Schedule = d.Get(commandResourceScheduleFieldName).(string)
	}

	if res, err := client.CommandClient.Update(command); err != nil {
		log.Println("Failed to update command")
		return fmt.Errorf("failed to update command: %w", err)
	} else {
		d.SetId(strconv.Itoa(res.ID))
	}
	return resourceCommandRead(d, m)
}

func resourceCommandDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	cid, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	if err := client.CommandClient.Delete(cid); err != nil {
		return err
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

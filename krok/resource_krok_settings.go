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
	commandSettingsKeyFieldName       = "key"
	commandSettingsValueFieldName     = "value"
	commandSettingsCommandIDFieldName = "command_id"
	commandSettingsInVaultFieldName   = "in_vault"
)

func resourceCommandSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceCommandSettingsCreate,
		Read:   resourceCommandSettingRead,
		Update: resourceCommandSettingUpdate,
		Delete: resourceCommandSettingDelete,

		Schema: map[string]*schema.Schema{
			commandSettingsKeyFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			commandSettingsValueFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			commandSettingsCommandIDFieldName: {
				Type:     schema.TypeInt,
				Required: true,
			},
			commandSettingsInVaultFieldName: {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

// resourceCommandSettingsCreate creates a Krok platform.
func resourceCommandSettingsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	expandedSetting, err := expandCommandSettingResource(d)
	if err != nil {
		return err
	}
	setting, err := client.SettingsClient.Create(expandedSetting)
	if err != nil {
		log.Println("Failed to create setting.")
		return fmt.Errorf("failed to create setting: %w", err)
	}
	d.SetId(strconv.Itoa(setting.ID))
	return resourceCommandRead(d, m)
}

// expandCommandSettingResource creates a Krok command setting structure out of a Terraform schema model.
func expandCommandSettingResource(d *schema.ResourceData) (*models.CommandSetting, error) {
	var (
		key, value string
		inVault    bool
		commandID  int
	)
	if v, ok := d.GetOk(commandSettingsKeyFieldName); ok {
		key = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", commandSettingsKeyFieldName)
	}
	if v, ok := d.GetOk(commandSettingsValueFieldName); ok {
		value = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", commandSettingsValueFieldName)
	}
	if v, ok := d.GetOk(commandSettingsInVaultFieldName); ok {
		inVault = v.(bool)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", commandSettingsInVaultFieldName)
	}
	if v, ok := d.GetOk(commandSettingsCommandIDFieldName); ok {
		commandID = v.(int)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", commandSettingsCommandIDFieldName)
	}
	return &models.CommandSetting{
		CommandID: commandID,
		Key:       key,
		Value:     value,
		InVault:   inVault,
	}, nil
}

// resourceCommandSettingUpdate updates command setting information from terraform stores.
func resourceCommandSettingUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	cid, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	setting, err := client.SettingsClient.Get(cid)
	if err != nil {
		log.Println("Failed to find command setting")
		d.SetId("")
		return err
	}

	if d.HasChange(commandSettingsValueFieldName) {
		setting.Value = d.Get(commandSettingsValueFieldName).(string)
	}

	if err := client.SettingsClient.Update(setting); err != nil {
		log.Println("Failed to update command setting")
		return fmt.Errorf("failed to update command setting: %w", err)
	} else {
		d.SetId(strconv.Itoa(setting.ID))
	}
	return resourceCommandSettingRead(d, m)
}

// resourceCommandSettingRead retrieves command setting information from terraform stores.
func resourceCommandSettingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	sid, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	setting, err := client.SettingsClient.Get(sid)
	if err != nil {
		log.Println("Failed to find command setting")
		d.SetId("")
		return err
	}

	for k, v := range flattenCommandSetting(setting) {
		if err := d.Set(k, v); err != nil {
			d.SetId("")
			return err
		}
	}
	return nil
}

// flattenCommandSetting creates a map from a command setting for easy storage on terraform.
func flattenCommandSetting(setting *models.CommandSetting) map[string]interface{} {
	flatCommandSetting := map[string]interface{}{
		commandSettingsKeyFieldName:       setting.Key,
		commandSettingsValueFieldName:     setting.Value,
		commandSettingsCommandIDFieldName: setting.CommandID,
		commandSettingsInVaultFieldName:   setting.InVault,
	}
	return flatCommandSetting
}

func resourceCommandSettingDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	sid, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	if err := client.SettingsClient.Delete(sid); err != nil {
		return err
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

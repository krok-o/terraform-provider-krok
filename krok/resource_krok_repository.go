package krok

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/krok-o/krok/pkg/models"

	"github.com/krok-o/terraform-provider-krok/pkg"
)

const (
	repoNameFieldName            = "name"
	repoURLFieldName             = "url"
	repoVCSFieldName             = "vcs"
	repoGitlabFieldName          = "gitlab"
	repoGitlabProjectIDFieldName = "project_id"
	repoAuthFieldName            = "auth"
	repoAuthSecretFieldName      = "secret"
	repoCommandsFieldName        = "commands"
	// repoEventsFieldName                      = "name"
)

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryCreate,
		Read:   resourceRepositoryRead,
		Update: resourceRepositoryUpdate,
		Delete: resourceRepositoryDelete,

		Schema: map[string]*schema.Schema{
			repoNameFieldName: {
				Type:     schema.TypeString,
				Required: true,
			},
			repoCommandsFieldName: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			repoAuthFieldName: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "1" && new == "0"
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						repoAuthSecretFieldName: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			repoGitlabFieldName: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "1" && new == "0"
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						repoGitlabProjectIDFieldName: {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

// resourceRepositoryCreate creates a Krok repository.
func resourceRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	expandedRepo, err := expandRepositoryResource(client, d)
	if err != nil {
		return err
	}
	if repo, err := client.RepositoryClient.Create(expandedRepo); err != nil {
		log.Println("Failed to create repository.")
		return err
	} else {
		d.SetId(strconv.Itoa(repo.ID))
	}
	return resourceRepositoryRead(d, m)
}

// authFields is a convenient wrapper around the auth schema for easy parsing
type authFields struct {
	secret string
}

// gitlabFields is a convenient wrapper around the gitlab schema for easy parsing
type gitlabFields struct {
	projectID int
}

// expandRepositoryResource creates a Krok repository structure out of a terraform schema model.
func expandRepositoryResource(client *pkg.KrokClient, d *schema.ResourceData) (*models.Repository, error) {
	var (
		name   string
		url    string
		vcs    int
		gitlab gitlabFields
		auth   authFields
		err    error
	)
	if v, ok := d.GetOk(repoNameFieldName); ok {
		name = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", repoNameFieldName)
	}
	if v, ok := d.GetOk(repoURLFieldName); ok {
		url = v.(string)
	} else {
		return nil, fmt.Errorf("unable to find parse field %s", repoURLFieldName)
	}
	if v, ok := d.GetOk(repoAuthFieldName); ok {
		if auth, err = expandAuth(v.([]interface{})); err != nil {
			return nil, err
		}
	}
	repo := &models.Repository{
		Name: name,
		URL:  url,
		VCS:  vcs,
		Auth: &models.Auth{
			Secret: auth.secret,
		},
	}
	if v, ok := d.GetOk(repoCommandsFieldName); ok {
		commands, err := expandCommands(client, v.([]int))
		if err != nil {
			return nil, err
		}
		repo.Commands = commands
	}
	if v, ok := d.GetOk(repoGitlabFieldName); ok {
		if gitlab, err = expandGitlab(v.([]interface{})); err != nil {
			return nil, err
		}
		repo.GitLab = &models.GitLab{
			ProjectID: gitlab.projectID,
		}
	}
	return repo, nil
}

// expandAuth gathers security data from the terraform store
func expandAuth(s []interface{}) (auth authFields, err error) {
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[repoAuthSecretFieldName]; ok {
			auth.secret = i.(string)
		} else {
			return auth, errors.New("secret field not found in auth")
		}
	}
	return
}

// expandAuth gathers security data from the terraform store
func expandGitlab(s []interface{}) (gitlab gitlabFields, err error) {
	for _, v := range s {
		item := v.(map[string]interface{})
		if i, ok := item[repoGitlabProjectIDFieldName]; ok {
			gitlab.projectID = i.(int)
		} else {
			return gitlab, errors.New("if gitlab is defined, project ID must be provided")
		}
	}
	return
}

// expandCommands gathers all commands for which the IDs have been defined.
func expandCommands(client *pkg.KrokClient, s []int) (commands []*models.Command, err error) {
	for _, v := range s {
		command, err := client.CommandClient.Get(v)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve command with id %d with error: %w", v, err)
		}
		commands = append(commands, command)
	}
	return
}

// resourceRepositoryRead retrieves repository information from terraform stores.
func resourceRepositoryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	rid, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	repo, err := client.RepositoryClient.Get(rid)
	if err != nil {
		log.Println("Failed to find repository")
		d.SetId("")
		return err
	}

	for k, v := range flattenRepository(repo) {
		if err := d.Set(k, v); err != nil {
			d.SetId("")
			return err
		}
	}
	return nil
}

// flattenRepository creates a map from a repository for easy storage on terraform.
func flattenRepository(repo *models.Repository) map[string]interface{} {
	commands := make([]int, 0)
	for _, c := range repo.Commands {
		commands = append(commands, c.ID)
	}
	flatRepo := map[string]interface{}{
		repoNameFieldName:     repo.Name,
		repoURLFieldName:      repo.URL,
		repoVCSFieldName:      repo.VCS,
		repoAuthFieldName:     flattenAuth(repo),
		repoCommandsFieldName: commands,
	}
	if repo.GitLab != nil {
		flatRepo[repoGitlabFieldName] = flattenGitlab(repo)
	}
	return map[string]interface{}{}
}

// flattenAuth takes the auth part of a repositroy and creates a sub map for terraform schema.
func flattenAuth(repo *models.Repository) []interface{} {
	return []interface{}{
		map[string]interface{}{
			repoAuthSecretFieldName: repo.Auth.Secret,
		},
	}
}

// flattenGitlab takes the gitlab part of a repositroy and creates a sub map for terraform schema.
func flattenGitlab(repo *models.Repository) []interface{} {
	return []interface{}{
		map[string]interface{}{
			repoGitlabProjectIDFieldName: repo.GitLab.GetProjectID(),
		},
	}
}

// resourceRepositoryUpdate checks fields for differences and updates a repository if necessary.
func resourceRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	rid, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	repo, err := client.RepositoryClient.Get(rid)
	if err != nil {
		log.Println("Failed to find repository")
		d.SetId("")
		return err
	}

	if d.HasChange(repoNameFieldName) {
		repo.Name = d.Get(repoNameFieldName).(string)
	}

	if res, err := client.RepositoryClient.Update(repo); err != nil {
		log.Println("Failed to update repositroy")
		return err
	} else {
		d.SetId(strconv.Itoa(res.ID))
	}
	return resourceRepositoryRead(d, m)
}

func resourceRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*pkg.KrokClient)
	rid, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	if err := client.RepositoryClient.Delete(rid); err != nil {
		return err
	}
	d.SetId("") // called automatically, but added to be explicit
	return nil
}

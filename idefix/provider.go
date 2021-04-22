package idefix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linkbynet/goidefix"
	"github.com/linkbynet/goidefix/services/authentification"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IDEFIX_URL", ""),
				Description: "This can be used to override the base URL for Idefix API. This can also be sourced from the `IDEFIX_URL` environment variable.",
			},
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IDEFIX_LOGIN", ""),
				Description: "The login wich should be used. This can also be sourced from the `IDEFIX_LOGIN` environment variable.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("IDEFIX_PASSWORD", ""),
				Description: "The password wich should be used. This can also be sourced from the `IDEFIX_PASSWORD` environment variable.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"idefix_project": resourceProject(),
			"idefix_ci":      resourceCI(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"idefix_project":  dataSourceProject(),
			"idefix_projects": dataSourceProjects(),
			"idefix_ci":       dataSourceCI(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var client *goidefix.Idefix
	var err error

	url := d.Get("url").(string)
	login := d.Get("login").(string)
	password := d.Get("password").(string)

	if url == "" {
		client, err = goidefix.New(ctx)
	} else {
		client, err = goidefix.NewWithEndpoint(ctx, url)
	}
	if err != nil {
		return nil, diag.FromErr(err)
	}
	_, err = client.Authentification.Login(ctx, &authentification.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, diags
}

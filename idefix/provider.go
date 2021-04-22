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
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IDEFIX_URL", "https://extranet.linkbynet.com/api/1.1"),
			},
			"login": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IDEFIX_LOGIN", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("IDEFIX_PASSWORD", nil),
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

	url := d.Get("url").(string)
	login := d.Get("login").(string)
	password := d.Get("password").(string)

	client, err := goidefix.NewWithEndpoint(ctx, url)
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

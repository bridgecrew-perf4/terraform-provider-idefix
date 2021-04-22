package idefix

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linkbynet/goidefix"
	"github.com/linkbynet/goidefix/services/project"
)

func dataSourceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectsRead,
		Description: "Use this data source to access information about existing Projects.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of this resource.",
			},
			"name_filter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name to filter the list of projects.",
			},
			"projects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The projects list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the project.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Name of the project.",
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*goidefix.Idefix)
	resp, err := client.Project.Search(ctx, &project.SearchRequest{
		Name: d.Get("name_filter").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	projects := flattenProjectsData(resp)
	if err := d.Set("projects", projects); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenProjectsData(projects *[]project.SearchResponse) []interface{} {
	if projects != nil {
		ps := make([]interface{}, len(*projects))

		for i, project := range *projects {
			p := make(map[string]interface{})

			p["id"] = project.ID
			p["name"] = project.Name

			ps[i] = p
		}

		return ps
	}

	return make([]interface{}, 0)
}

package idefix

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linkbynet/goidefix"
	"github.com/linkbynet/goidefix/services/project"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Description: "Use this data source to access information about an existing Project.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of this resource.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the project.",
			},
			"company_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The company ID.",
			},
			"parent_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the parent project.",
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := strconv.Itoa(d.Get("id").(int))

	client := m.(*goidefix.Idefix)
	project, err := client.Project.Read(ctx, &project.ReadRequest{
		ID: id,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", project.Name)
	d.Set("company_id", project.CompanyID)
	d.Set("parent_id", project.ParentID)

	d.SetId(id)

	return diags
}

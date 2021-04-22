package idefix

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linkbynet/goidefix"
	"github.com/linkbynet/goidefix/services/ci"
)

func dataSourceCI() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCIRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"type_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"company_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"project_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"outsourcing_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_level_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"team": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_owner_lbn": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCIRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := strconv.Itoa(d.Get("id").(int))

	client := m.(*goidefix.Idefix)
	ci, err := client.CI.Read(ctx, &ci.ReadRequest{
		ID: id,
	})
	if err != nil {
		diag.FromErr(err)
	}

	var projectIDs []int
	pids := strings.Split(ci.ProjectIDs, ",")
	for _, pid := range pids {
		if pid != "" {
			id, err := strconv.Atoi(pid)
			if err != nil {
				return diag.FromErr(err)
			}

			projectIDs = append(projectIDs, id)
		}
	}

	d.Set("name", ci.Name)
	d.Set("company_id", ci.CompanyID)
	d.Set("type_id", ci.TypeID)
	d.Set("company_id", ci.CompanyID)
	d.Set("project_ids", projectIDs)
	d.Set("outsourcing_name", ci.OutSourcingName)
	d.Set("service_level_id", ci.ServiceLevelID)
	d.Set("team", ci.Team)
	d.Set("is_owner_lbn", ci.IsOwnerLBN)
	d.Set("comment", ci.Comment)

	d.SetId(id)

	return diags
}

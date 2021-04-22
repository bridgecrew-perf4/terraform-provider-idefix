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

func resourceCI() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCICreate,
		ReadContext:   resourceCIRead,
		UpdateContext: resourceCIUpdate,
		DeleteContext: resourceCIDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  41,
			},
			"company_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"project_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"outsourcing_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0 - Non-défini",
			},
			"service_level_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100000080,
			},
			"team": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Unix",
			},
			"is_owner_lbn": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCICreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goidefix.Idefix)

	ids := d.Get("project_ids").([]interface{})
	projectIDs := make([]int, len(ids))
	for i := range ids {
		projectIDs[i] = ids[i].(int)
	}

	_, err := client.CI.Create(ctx, &ci.CreateRequest{
		Name:            d.Get("name").(string),
		TypeID:          d.Get("type_id").(int),
		CompanyID:       d.Get("company_id").(int),
		ProjectIDs:      projectIDs,
		OutSourcingName: d.Get("outsourcing_name").(string),
		ServiceLevelID:  d.Get("service_level_id").(int),
		Team:            d.Get("team").(string),
		IsOwnerLBN:      d.Get("is_owner_lbn").(bool),
		Comment:         d.Get("comment").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCIRead(ctx, d, m)
}

func resourceCIRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*goidefix.Idefix)
	ci, err := client.CI.Read(ctx, &ci.ReadRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
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

	d.SetId(d.Id())
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

	return diags
}

func resourceCIUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goidefix.Idefix)

	ids := d.Get("project_ids").([]interface{})
	projectIDs := make([]int, len(ids))
	for i := range ids {
		projectIDs[i] = ids[i].(int)
	}

	_, err := client.CI.Update(ctx, &ci.UpdateRequest{
		ID:              d.Id(),
		Name:            d.Get("name").(string),
		TypeID:          d.Get("type_id").(int),
		CompanyID:       d.Get("company_id").(int),
		ProjectIDs:      projectIDs,
		OutSourcingName: d.Get("outsourcing_name").(string),
		ServiceLevelID:  d.Get("service_level_id").(int),
		Team:            d.Get("team").(string),
		IsOwnerLBN:      d.Get("is_owner_lbn").(bool),
		Comment:         d.Get("comment").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCIRead(ctx, d, m)
}

func resourceCIDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*goidefix.Idefix)
	_, err := client.CI.Delete(ctx, &ci.DeleteRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

package idefix

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linkbynet/goidefix"
	"github.com/linkbynet/goidefix/services/ci"
	"github.com/linkbynet/goidefix/services/equipment"
	"github.com/linkbynet/goidefix/services/monitoring"
)

func resourceCI() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCICreate,
		ReadContext:   resourceCIRead,
		UpdateContext: resourceCIUpdate,
		DeleteContext: resourceCIDelete,
		Description:   "Manages CI.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of this resource.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this CI.",
			},
			"type_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     41,
				Description: "The type of the CI.",
			},
			"company_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The company ID associated to the CI.",
			},
			"project_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The projects associated to the CI.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"outsourcing_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0 - Non-défini",
				Description: "The Outsourcing level name.",
			},
			"service_level_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100000080,
				Description: "The Level of the service.",
			},
			"team": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Unix",
				Description: "The team in charge.",
			},
			"is_owner_lbn": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The owner of the CI.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

	ci, err := client.CI.Create(ctx, &ci.CreateRequest{
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

	d.SetId(ci.ID)

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

	typeID, err := strconv.Atoi(ci.TypeID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Id())
	d.Set("name", ci.Name)
	d.Set("company_id", ci.CompanyID)
	d.Set("type_id", typeID)
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

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	events, err := client.Monitoring.SearchEvents(ctx, &monitoring.SearchEventsRequest{
		EquipmentIDs: []int{id},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, event := range *events {
		_, err := client.Monitoring.DeleteEvents(ctx, &monitoring.DeleteEventsRequest{
			ID: event.ID,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	_, err = client.Equipment.Delete(ctx, &equipment.DeleteRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

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
			"service_at": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"required_services": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"monitoring_tool": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"key_dates": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"function_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"service_cloud": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"product_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

	cir, err := client.CI.Create(ctx, &ci.CreateRequest{
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

	if v, ok := d.GetOk("service_cloud"); ok && v.(*schema.Set).Len() > 0 {
		for _, serviceCloudSet := range v.(*schema.Set).List() {
			var subscriptionId, productID int
			var regionID string

			serviceCloud, ok := serviceCloudSet.(map[string]interface{})

			if !ok {
				continue
			}

			if v, ok := serviceCloud["subscription_id"].(int); ok && v > 0 {
				subscriptionId = v
			}

			if v, ok := serviceCloud["product_id"].(int); ok && v > 0 {
				productID = v
			}

			if v, ok := serviceCloud["region_id"].(string); ok && v != "" {
				regionID = v
			}

			_, err := client.CI.UpdateServiceCloud(ctx, &ci.UpdateServiceCloudRequest{
				ID:             cir.ID,
				SubscriptionID: subscriptionId,
				ProductID:      productID,
				RegionID:       regionID,
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if v, ok := d.GetOk("key_dates"); ok && v.(*schema.Set).Len() > 0 {
		for _, keyDatesSet := range v.(*schema.Set).List() {
			var envIDs, funcIDs []int

			keyDates, ok := keyDatesSet.(map[string]interface{})

			if !ok {
				continue
			}

			if v, ok := keyDates["environment_ids"].([]interface{}); ok && len(v) > 0 {
				for i := range v {
					envIDs = append(envIDs, v[i].(int))
				}
			}

			if v, ok := keyDates["function_ids"].([]interface{}); ok && len(v) > 0 {
				for i := range v {
					funcIDs = append(funcIDs, v[i].(int))
				}
			}

			_, err := client.CI.UpdateUseAndKeyDate(ctx, &ci.UpdateUseAndKeyDateRequest{
				ID:             cir.ID,
				EnvSelect:      0,
				EnvironmentIDs: envIDs,
				FuncSelect:     0,
				FunctionIDs:    funcIDs,
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if v, ok := d.GetOk("service_at"); ok && v.(*schema.Set).Len() > 0 {
		for _, serviceATSet := range v.(*schema.Set).List() {
			var requiredServices, monitoringTool string

			serviceAT, ok := serviceATSet.(map[string]interface{})

			if !ok {
				continue
			}

			if v, ok := serviceAT["required_services"].([]interface{}); ok && len(v) > 0 {
				var ids []string
				for i := range v {
					ids = append(ids, strconv.Itoa(v[i].(int)))
				}
				requiredServices = strings.Join(ids, ",")
			}

			if v, ok := serviceAT["monitoring_tool"].([]interface{}); ok && len(v) > 0 {
				var ids []string
				for i := range v {
					ids = append(ids, strconv.Itoa(v[i].(int)))
				}
				monitoringTool = strings.Join(ids, ",")
			}

			_, err := client.Equipment.UpdateAT(ctx, &equipment.UpdateATRequest{
				ID:               cir.ID,
				RequiredServices: requiredServices,
				MonitoringTool:   monitoringTool,
				BackupComment:    "Asset PaaS",
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	_, err = client.CI.UpdatePlatform(ctx, &ci.UpdatePlatformRequest{
		ID: cir.ID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cir.ID)

	return resourceCIRead(ctx, d, m)
}

func resourceCIRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*goidefix.Idefix)
	cir, err := client.CI.Read(ctx, &ci.ReadRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var projectIDs []int
	pids := strings.Split(cir.ProjectIDs, ",")
	for _, pid := range pids {
		if pid != "" {
			id, err := strconv.Atoi(pid)
			if err != nil {
				return diag.FromErr(err)
			}

			projectIDs = append(projectIDs, id)
		}
	}

	typeID, err := strconv.Atoi(cir.TypeID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Id())
	d.Set("name", cir.Name)
	d.Set("company_id", cir.CompanyID)
	d.Set("type_id", typeID)
	d.Set("company_id", cir.CompanyID)
	d.Set("project_ids", projectIDs)
	d.Set("outsourcing_name", cir.OutSourcingName)
	d.Set("service_level_id", cir.ServiceLevelID)
	d.Set("team", cir.Team)
	d.Set("is_owner_lbn", cir.IsOwnerLBN)
	d.Set("comment", cir.Comment)

	sc, err := client.CI.ReadServiceCloud(ctx, &ci.ReadServiceCloudRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	serviceCloud := map[string]interface{}{}
	serviceCloud["subscription_id"] = sc.SubscriptionID
	serviceCloud["product_id"] = sc.ProductID
	serviceCloud["region_id"] = sc.RegionID

	var serviceCloudSet []interface{}
	serviceCloudSet = append(serviceCloudSet, serviceCloud)
	d.Set("service_cloud", serviceCloudSet)

	kd, err := client.CI.ReadUseAndKeyDate(ctx, &ci.ReadUseAndKeyDateRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	keyDate := map[string]interface{}{}

	var envIDs []int
	envIDsList := strings.Split(kd.EnvironmentIDs, ",")
	if len(envIDsList) > 0 {
		for _, envID := range envIDsList {
			if envID == "" {
				continue
			}

			id, err := strconv.Atoi(envID)
			if err != nil {
				return diag.FromErr(err)
			}

			envIDs = append(envIDs, id)
		}
	}
	keyDate["environment_ids"] = envIDs

	var funcIDs []int
	funcIDsList := strings.Split(kd.FunctionIDs, ",")
	if len(funcIDsList) > 0 {
		for _, funcID := range funcIDsList {
			if funcID == "" {
				continue
			}

			id, err := strconv.Atoi(funcID)
			if err != nil {
				return diag.FromErr(err)
			}

			funcIDs = append(funcIDs, id)
		}
	}
	keyDate["function_ids"] = funcIDs

	var keyDatesSet []interface{}
	keyDatesSet = append(keyDatesSet, keyDate)
	d.Set("key_dates", keyDatesSet)

	at, err := client.Equipment.ReadAT(ctx, &equipment.ReadATRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	serviceAT := map[string]interface{}{}

	var requiredServices []int
	requiredServicesList := strings.Split(at.RequiredServices, ",")
	if len(requiredServicesList) > 0 {
		for _, requiredService := range requiredServicesList {
			if requiredService == "" {
				continue
			}

			id, err := strconv.Atoi(requiredService)
			if err != nil {
				return diag.FromErr(err)
			}

			requiredServices = append(requiredServices, id)
		}
	}
	serviceAT["required_services"] = requiredServices

	var monitoringTools []int
	monitoringToolList := strings.Split(at.MonitoringTool, ",")
	if len(monitoringToolList) > 0 {
		for _, monitoringTool := range monitoringToolList {
			if monitoringTool == "" {
				continue
			}

			id, err := strconv.Atoi(monitoringTool)
			if err != nil {
				return diag.FromErr(err)
			}

			monitoringTools = append(monitoringTools, id)
		}
	}
	serviceAT["monitoring_tool"] = monitoringTools

	var serviceATSet []interface{}
	serviceATSet = append(serviceATSet, serviceAT)
	d.Set("service_at", serviceATSet)

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

	if v, ok := d.GetOk("service_cloud"); ok && v.(*schema.Set).Len() > 0 {
		for _, serviceCloudSet := range v.(*schema.Set).List() {
			var subscriptionId, productID int
			var regionID string

			serviceCloud, ok := serviceCloudSet.(map[string]interface{})

			if !ok {
				continue
			}

			if v, ok := serviceCloud["subscription_id"].(int); ok && v > 0 {
				subscriptionId = v
			}

			if v, ok := serviceCloud["product_id"].(int); ok && v > 0 {
				productID = v
			}

			if v, ok := serviceCloud["region_id"].(string); ok && v != "" {
				regionID = v
			}

			_, err := client.CI.UpdateServiceCloud(ctx, &ci.UpdateServiceCloudRequest{
				ID:             d.Id(),
				SubscriptionID: subscriptionId,
				ProductID:      productID,
				RegionID:       regionID,
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if v, ok := d.GetOk("key_dates"); ok && v.(*schema.Set).Len() > 0 {
		for _, keyDatesSet := range v.(*schema.Set).List() {
			var envIDs, funcIDs []int

			keyDates, ok := keyDatesSet.(map[string]interface{})

			if !ok {
				continue
			}

			if v, ok := keyDates["environment_ids"].([]interface{}); ok && len(v) > 0 {
				for i := range v {
					envIDs = append(envIDs, v[i].(int))
				}
			}

			if v, ok := keyDates["function_ids"].([]interface{}); ok && len(v) > 0 {
				for i := range v {
					funcIDs = append(funcIDs, v[i].(int))
				}
			}

			_, err := client.CI.UpdateUseAndKeyDate(ctx, &ci.UpdateUseAndKeyDateRequest{
				ID:             d.Id(),
				EnvSelect:      0,
				EnvironmentIDs: envIDs,
				FuncSelect:     0,
				FunctionIDs:    funcIDs,
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if v, ok := d.GetOk("service_at"); ok && v.(*schema.Set).Len() > 0 {
		for _, serviceATSet := range v.(*schema.Set).List() {
			var requiredServices, monitoringTool string

			serviceAT, ok := serviceATSet.(map[string]interface{})

			if !ok {
				continue
			}

			if v, ok := serviceAT["required_services"].([]interface{}); ok && len(v) > 0 {
				var ids []string
				for i := range v {
					ids = append(ids, strconv.Itoa(v[i].(int)))
				}
				requiredServices = strings.Join(ids, ",")
			}

			if v, ok := serviceAT["monitoring_tool"].([]interface{}); ok && len(v) > 0 {
				var ids []string
				for i := range v {
					ids = append(ids, strconv.Itoa(v[i].(int)))
				}
				monitoringTool = strings.Join(ids, ",")
			}

			_, err := client.Equipment.UpdateAT(ctx, &equipment.UpdateATRequest{
				ID:               d.Id(),
				RequiredServices: requiredServices,
				MonitoringTool:   monitoringTool,
				BackupComment:    "Asset PaaS",
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	_, err = client.CI.UpdatePlatform(ctx, &ci.UpdatePlatformRequest{
		ID: d.Id(),
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

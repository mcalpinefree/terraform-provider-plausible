package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mcalpinefree/terraform-provider-plausible/plausibleclient"
)

func resourceGoal() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoalCreate,
		Read:   resourceGoalRead,
		Delete: resourceGoalDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The goal ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site_id": {
				Description: "The domain of the site to create the goal for.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"page_path": {
				Description:   "Page path event. E.g. `/success`",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"event_name"},
			},
			"event_name": {
				Description:   "Custom event E.g. `Signup`",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"page_path"},
			},
		},
	}
}

func resourceGoalCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient)
	domain := d.Get("site_id").(string)

	var goalType plausibleclient.GoalType
	goal := ""
	if v, ok := d.GetOk("page_path"); ok {
		goalType = plausibleclient.PagePath
		goal = v.(string)
	} else if v, ok := d.GetOk("event_name"); ok {
		goalType = plausibleclient.EventName
		goal = v.(string)
	} else {
		return fmt.Errorf("page_path or event_name needs to be defined")
	}

	resp, err := client.plausibleClient.CreateGoal(domain, goalType, goal)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", resp.ID))

	return resourceGoalSetResourceData(resp, d)
}

func resourceGoalSetResourceData(g *plausibleclient.Goal, d *schema.ResourceData) error {
	d.Set("site_id", g.Domain)
	if g.PagePath != nil {
		d.Set("page_path", *g.PagePath)
	} else if g.EventName != nil {
		d.Set("event_name", *g.EventName)
	} else {
		return fmt.Errorf("either PagePath or EventName needs to not be nil")
	}
	return nil
}

func resourceGoalRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	g := &plausibleclient.Goal{
		ID:     idInt,
		Domain: d.Get("site_id").(string),
	}

	if v, ok := d.GetOk("page_path"); ok {
		pagePath := v.(string)
		g.PagePath = &pagePath
	} else if v, ok := d.GetOk("event_name"); ok {
		eventName := v.(string)
		g.EventName = &eventName
	} else {
		return fmt.Errorf("page_path or event_name needs to be defined")
	}

	return resourceGoalSetResourceData(g, d)
}

func resourceGoalDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient)
	id := d.Id()
	domain := d.Get("site_id").(string)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return client.plausibleClient.DeleteGoal(domain, idInt)
}

package extras

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"awsextras_terminate_instances": resourceTerminateEc2Instances(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(d.Get("region").(string)))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return cfg, nil
}

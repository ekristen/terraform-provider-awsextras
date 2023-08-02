package extras

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecretsManagerSecretRotation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Secret ID",
				Required:    true,
			},
			"automatically_after_days": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "rotate the secret every N days",
				Optional:    true,
			},
		},
		ReadContext:   resourceSecretsManagerSecretRotationRead,
		CreateContext: resourceSecretsManagerSecretRotationCreate,
		UpdateContext: resourceSecretsManagerSecretRotationUpdate,
		DeleteContext: resourceSecretsManagerSecretRotationDelete,
	}
}

func resourceSecretsManagerSecretRotationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(aws.Config)
	client := secretsmanager.NewFromConfig(cfg)

	secretId := d.Get("secret_id").(string)

	s, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretId),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("automatically_after_days", aws.ToInt64(s.RotationRules.AutomaticallyAfterDays)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSecretsManagerSecretRotationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(aws.Config)
	client := secretsmanager.NewFromConfig(cfg)

	secretId := d.Get("secret_id").(string)
	automaticallyAfterDays := d.Get("automatically_after_days").(int)

	s, err := client.RotateSecret(ctx, &secretsmanager.RotateSecretInput{
		SecretId: aws.String(secretId),
		RotationRules: &types.RotationRulesType{
			AutomaticallyAfterDays: aws.Int64(int64(automaticallyAfterDays)),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aws.ToString(s.ARN))
	if err := d.Set("automatically_after_days", automaticallyAfterDays); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSecretsManagerSecretRotationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(aws.Config)
	client := secretsmanager.NewFromConfig(cfg)

	secretId := d.Get("secret_id").(string)
	automaticallyAfterDays := d.Get("automatically_after_days").(int)

	s, err := client.RotateSecret(ctx, &secretsmanager.RotateSecretInput{
		SecretId: aws.String(secretId),
		RotationRules: &types.RotationRulesType{
			AutomaticallyAfterDays: aws.Int64(int64(automaticallyAfterDays)),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aws.ToString(s.ARN))
	if err := d.Set("automatically_after_days", automaticallyAfterDays); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSecretsManagerSecretRotationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

package extras

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

func resourceRemoveKeyPairs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"include_regex": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "include key pair by regex matching on key pair name",
				Optional:    true,
			},
			"exclude_regex": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "exclude key pair by regex on key pair name",
				Optional:    true,
			},
			"exclude_key_pair_names": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "exclude any of these key pair names",
				Optional:    true,
			},
			"key_pair_names": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of computed ssh key pair names remove",
				Computed:    true,
				Optional:    false,
				Required:    false,
			},
		},
		ReadContext:   resourceRemoveKeyPairsRead,
		CreateContext: resourceRemoveKeyPairsCreate,
		UpdateContext: resourceRemoveKeyPairsUpdate,
		DeleteContext: resourceRemoveKeyPairsDelete,
	}
}

func resourceRemoveKeyPairsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(aws.Config)
	client := ec2.NewFromConfig(cfg)

	includeRegexSet := d.Get("include_regex").(*schema.Set)
	includeRegex := make([]string, includeRegexSet.Len())
	for i, regex := range includeRegexSet.List() {
		includeRegex[i] = regex.(string)
	}

	excludeRegexSet := d.Get("exclude_regex").(*schema.Set)
	excludeRegex := make([]string, excludeRegexSet.Len())
	for i, regex := range excludeRegexSet.List() {
		excludeRegex[i] = regex.(string)
	}

	excludeNamesSet := d.Get("exclude_key_pair_names").(*schema.Set)
	excludeNames := make([]string, excludeNamesSet.Len())
	for i, id := range excludeNamesSet.List() {
		excludeNames[i] = id.(string)
	}

	output, err := client.DescribeKeyPairs(ctx, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var keypairNames []string = make([]string, 0)

	for _, keypair := range output.KeyPairs {
		if slices.Contains(excludeNames, aws.ToString(keypair.KeyName)) {
			log.Printf("[DEBUG] excluding keypair: %s/%s, reason: explicit exclude", aws.ToString(keypair.KeyPairId), aws.ToString(keypair.KeyName))
			continue
		}

		for _, regex := range excludeRegex {
			re, err := regexp.Compile(regex)
			if err != nil {
				return diag.FromErr(fmt.Errorf("unable to compile regex: %s", err.Error()))
			}

			if re.MatchString(aws.ToString(keypair.KeyName)) {
				log.Printf("[DEBUG] excluding keypair: %s/%s, reason: exclude regex matched", aws.ToString(keypair.KeyPairId), aws.ToString(keypair.KeyName))
				continue
			}
		}

		for _, regex := range includeRegex {
			re, err := regexp.Compile(regex)
			if err != nil {
				return diag.FromErr(fmt.Errorf("unable to compile regex: %s", err.Error()))
			}

			if re.MatchString(aws.ToString(keypair.KeyName)) {
				log.Printf("[DEBUG] including keypair: %s/%s", aws.ToString(keypair.KeyPairId), aws.ToString(keypair.KeyName))
				keypairNames = append(keypairNames, aws.ToString(keypair.KeyName))
			} else {
				log.Printf("[DEBUG] excluding keypair: %s/%s, reason: include regex mismatched", aws.ToString(keypair.KeyPairId), aws.ToString(keypair.KeyName))
				continue
			}
		}
	}

	d.Set("key_pair_names", keypairNames)

	return nil
}

func resourceRemoveKeyPairsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return nil
}

func resourceRemoveKeyPairsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] UPDATE CALLED")
	return nil
}

func resourceRemoveKeyPairsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(aws.Config)
	client := ec2.NewFromConfig(cfg)

	keypairNamesSet := d.Get("key_pair_names").(*schema.Set)
	keypairNames := make([]string, keypairNamesSet.Len())
	for i, keypairName := range keypairNamesSet.List() {
		keypairNames[i] = keypairName.(string)
		log.Printf("[INFO] removing keypair: %s", keypairNames[i])
	}

	if len(keypairNames) == 0 {
		return nil
	}

	for _, id := range keypairNames {
		deleteInput := &ec2.DeleteKeyPairInput{
			KeyName: aws.String(id),
		}

		_, err := client.DeleteKeyPair(ctx, deleteInput)

		if err != nil {
			return diag.FromErr(fmt.Errorf("unable to remove keypair (%s): %w", id, err))
		}
	}

	return nil
}

package extras

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
	"log"
	"time"
)

const (
	ErrCodeInvalidInstanceIDNotFound = "InvalidInstanceID.NotFound"
)

func resourceTerminateEc2Instances() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"subnet_ids": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "known list of subnet ids to look for Terminated ec2 instances",
				Required:    true,
			},
			"include_untagged": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "considered untagged as Terminated",
				Default:     true,
				Optional:    true,
			},
			"exclude_instance_ids": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "known list of instance ids to treat as managed",
				Optional:    true,
			},
			"exclude_tags": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "consider resources with tags as managed (ie kubernetes)",
				Optional:    true,
			},
			"instance_ids": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of computed instance ids to remove",
				Computed:    true,
				Optional:    false,
				Required:    false,
			},
		},
		ReadContext:   resourceTerminateEc2InstancesRead,
		CreateContext: resourceTerminateEc2InstancesCreate,
		UpdateContext: resourceTerminateEc2InstancesUpdate,
		DeleteContext: resourceTerminateEc2InstancesDelete,
	}
}

func resourceTerminateEc2InstancesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(aws.Config)
	client := ec2.NewFromConfig(cfg)

	subnetIdsSet := d.Get("subnet_ids").(*schema.Set)
	subnetIds := make([]string, subnetIdsSet.Len())
	for i, subnetId := range subnetIdsSet.List() {
		subnetIds[i] = subnetId.(string)
	}

	excludeInstanceIdsSet := d.Get("exclude_instance_ids").(*schema.Set)
	excludeInstanceIds := make([]string, excludeInstanceIdsSet.Len())
	for i, instanceId := range excludeInstanceIdsSet.List() {
		excludeInstanceIds[i] = instanceId.(string)
	}

	excludeTagsSet := d.Get("exclude_tags").(*schema.Set)
	excludeTags := make([]string, excludeTagsSet.Len())
	for i, tag := range excludeTagsSet.List() {
		excludeTags[i] = tag.(string)
	}

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("subnet-id"),
				Values: subnetIds,
			},
		},
	}

	output, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	var instanceIds []string = make([]string, 0)

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			tid := aws.ToString(instance.InstanceId)
			log.Printf("[DEBUG] found instance: %s", tid)
			if slices.Contains(excludeInstanceIds, tid) {
				log.Printf("[DEBUG] skipping instance: %s, reason: excluded", tid)
				continue
			}

			if len(instance.Tags) != 1 {
				log.Printf("[DEBUG] skipping instance: %s, reason: more than 1 tag", tid)
				continue
			} else {
				for _, t := range instance.Tags {
					key := aws.ToString(t.Key)
					if slices.Contains(excludeTags, key) {
						continue
					}
				}
			}

			instanceIds = append(instanceIds, tid)
		}
	}

	d.Set("instance_ids", instanceIds)

	return nil
}

func resourceTerminateEc2InstancesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return nil
}

func resourceTerminateEc2InstancesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] UPDATE CALLED")
	return nil
}

func resourceTerminateEc2InstancesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(aws.Config)
	client := ec2.NewFromConfig(cfg)

	instanceIdsSet := d.Get("instance_ids").(*schema.Set)
	instanceIds := make([]string, instanceIdsSet.Len())
	for i, instanceId := range instanceIdsSet.List() {
		instanceIds[i] = instanceId.(string)
		log.Printf("[INFO] terminating ec2 instance: %s", instanceIds[i])
	}

	if len(instanceIds) == 0 {
		return nil
	}

	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIds,
	}

	_, err := client.TerminateInstances(ctx, terminateInput)

	if tfawserr.ErrCodeEquals(err, ErrCodeInvalidInstanceIDNotFound) {
		return nil
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("terminating EC2 Instance (%s): %w", instanceIds, err))
	}

	for _, id := range instanceIds {
		if _, err := WaitInstanceDeleted(ctx, client, id, 5*time.Minute); err != nil {
			return diag.FromErr(fmt.Errorf("waiting for EC2 Instance (%s) delete: %w", id, err))
		}
	}

	return nil
}

func WaitInstanceDeleted(ctx context.Context, client *ec2.Client, id string, timeout time.Duration) (*types.Instance, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(types.InstanceStateNamePending),
			string(types.InstanceStateNameRunning),
			string(types.InstanceStateNameShuttingDown),
			string(types.InstanceStateNameStopping),
			string(types.InstanceStateNameStopped),
		},
		Target:     []string{string(types.InstanceStateNameTerminated)},
		Refresh:    StatusInstanceState(ctx, client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*types.Instance); ok {
		if stateReason := output.StateReason; stateReason != nil {
			return output, fmt.Errorf("%s", aws.ToString(stateReason.Message))
		}

		return output, err
	}

	return nil, err
}

func StatusInstanceState(ctx context.Context, client *ec2.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Don't call FindInstanceByID as it maps useful status codes to NotFoundError.
		output, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
			InstanceIds: []string{id},
		})

		if err != nil {
			return nil, "", err
		}

		instances := output.Reservations[0].Instances

		if len(instances) == 0 || instances[0].State == nil {
			return nil, "", fmt.Errorf("empty result")
		}

		if count := len(instances); count > 1 {
			return nil, "", fmt.Errorf("too many results")
		}

		return instances[0], string(instances[0].State.Name), nil
	}
}

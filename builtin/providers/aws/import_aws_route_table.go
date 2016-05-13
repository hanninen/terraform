package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/schema"
)

// Route table import also imports all the rules
func resourceAwsRouteTableImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*AWSClient).ec2conn

	// First query the resource itself
	id := d.Id()
	resp, err := conn.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{&id},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.RouteTables) < 1 || resp.RouteTables[0] == nil {
		return nil, fmt.Errorf("route table %s is not found", id)
	}
	table := resp.RouteTables[0]

	// Start building our results
	results := make([]*schema.ResourceData, 1, 1+len(table.Routes))
	results[0] = d

	// Construct the routes
	subResource := resourceAwsRoute()
	for _, route := range table.Routes {
		// Minimal data for route
		d := subResource.Data(nil)
		d.SetType("aws_route")
		d.Set("route_table_id", id)
		d.Set("destination_cidr_block", route.DestinationCidrBlock)
		d.SetId(routeIDHash(d, route))
		results = append(results, d)
	}

	return results, nil
}

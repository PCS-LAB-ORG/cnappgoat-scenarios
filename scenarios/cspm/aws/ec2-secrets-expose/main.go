package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a new VPC
		vpc, err := ec2.NewVpc(ctx, "myvpc", &ec2.VpcArgs{
			CidrBlock: pulumi.String("10.0.0.0/16"),
		})
		if err != nil {
			return err
		}

		// Create a new subnet
		subnet, err := ec2.NewSubnet(ctx, "mysubnet", &ec2.SubnetArgs{
			VpcId:     vpc.ID(),
			CidrBlock: pulumi.String("10.0.1.0/24"),
		})
		if err != nil {
			return err
		}

		// Get the AMI
		mostRecent := true
		amiResult, err := aws.GetAmi(ctx, &aws.GetAmiArgs{
			Owners:     []string{"amazon"},
			MostRecent: &mostRecent,
			Filters:    []aws.GetAmiFilter{{Name: "name", Values: []string{"amzn2-ami-hvm-2.0.*-x86_64-ebs"}}},
		})
		if err != nil {
			return err
		}

		// Create a new EC2 instance
		_, err = ec2.NewInstance(ctx, "myec2instance", &ec2.InstanceArgs{
			InstanceType:             pulumi.String("t2.micro"),
			AssociatePublicIpAddress: pulumi.BoolPtr(true),
			SubnetId:                 subnet.ID(),
			Ami:                      pulumi.String(amiResult.Id),
			UserData: pulumi.String(`#!/bin/bash
                mysql -u cnaappgoat -p mysecretpassword1231`),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("my-ec2-instance131"),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}

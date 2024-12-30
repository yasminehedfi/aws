package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var region string

func createEC2Session() *ec2.EC2 {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Unable to create AWS session: %v", err)
	}
	return ec2.New(sess)
}

func listInstances() {
	svc := createEC2Session()

	result, err := svc.DescribeInstances(nil)
	if err != nil {
		log.Fatalf("Unable to describe instances: %v", err)
	}

	fmt.Println("Instances:")
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Printf("ID: %s, State: %s, Type: %s\n", *instance.InstanceId, *instance.State.Name, *instance.InstanceType)
		}
	}
}

func startInstance(instanceID string) {
	svc := createEC2Session()

	_, err := svc.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		log.Fatalf("Unable to start instance: %v", err)
	}

	fmt.Printf("Instance %s is starting...\n", instanceID)
}

func stopInstance(instanceID string) {
	svc := createEC2Session()

	_, err := svc.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		log.Fatalf("Unable to stop instance: %v", err)
	}

	fmt.Printf("Instance %s is stopping...\n", instanceID)
}
func createInstance(imageID, instanceType, keyName string) {
    svc := createEC2Session()

    runInput := &ec2.RunInstancesInput{
        ImageId:      aws.String(imageID),       
        InstanceType: aws.String(instanceType), 
        KeyName:      aws.String(keyName),      
        MinCount:     aws.Int64(1),             
        MaxCount:     aws.Int64(1),             
    }

    result, err := svc.RunInstances(runInput)
    if err != nil {
        log.Fatalf("Erreur lors de la création de l'instance : %v", err)
    }

    for _, instance := range result.Instances {
        fmt.Printf("Instance créée avec succès ! ID: %s, Type: %s, AMI: %s\n",
            *instance.InstanceId, *instance.InstanceType, *instance.ImageId)
    }
}


var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all EC2 instances",
	Run: func(cmd *cobra.Command, args []string) {
		listInstances()
	},
}

var startCmd = &cobra.Command{
	Use:   "start [instanceID]",
	Short: "Start EC2 instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startInstance(args[0])
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [instanceID]",
	Short: "Stop EC2 instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stopInstance(args[0])
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "AWS region")
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(createCmd)
}

var rootCmd = &cobra.Command{
	Use:   "ec2-cli",
	Short: "A CLI tool to interact with AWS EC2",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("EC2 CLI: Use `list`,‘create‘, `start`, or `stop` commands")
	},

}

var createCmd = &cobra.Command{
    Use:   "create [imageID] [instanceType] [keyName]",
    Short: "Créer une nouvelle instance EC2",
    Args:  cobra.ExactArgs(3), 
    Run: func(cmd *cobra.Command, args []string) {
        imageID := args[0]
        instanceType := args[1]
        keyName := args[2]
        createInstance(imageID, instanceType, keyName)
    },
}


func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

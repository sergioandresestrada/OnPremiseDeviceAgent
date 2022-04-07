package database

import (
	"backend/pkg/types"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDB defines the struct used to implement Database interface using AWS DynamoDB
// It contains an DynamoDB client and the name of the tables to be used
type DynamoDB struct {
	dynamoDBClient      *dynamodb.Client
	DEVICES_TABLE_NAME  string
	MESSAGES_TABLE_NAME string
}

// NewDatabaseDynamoDB creates and returns the reference to a new DynamoDB struct
func NewDatabaseDynamoDB() *DynamoDB {
	db := &DynamoDB{}
	db.initialize()
	return db
}

func (db *DynamoDB) initialize() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-3"))

	if err != nil {
		panic(fmt.Sprintf("Configuration error in AWS DynamoDB: %v\n", err))
	}

	_, ok := os.LookupEnv("DYNAMO_DB_DEVICES_TABLE_NAME")
	if !ok {
		panic("Environment variable DYNAMO_DB_DEVICES_TABLE_NAME does not exist")
	}

	_, ok = os.LookupEnv("DYNAMO_DB_MESSAGES_TABLE_NAME")
	if !ok {
		panic("Environment variable DYNAMO_DB_MESSAGES_TABLE_NAME does not exist")
	}

	db.DEVICES_TABLE_NAME = os.Getenv("DYNAMO_DB_DEVICES_TABLE_NAME")
	db.MESSAGES_TABLE_NAME = os.Getenv("DYNAMO_DB_MESSAGES_TABLE_NAME")

	db.dynamoDBClient = dynamodb.NewFromConfig(cfg)
}

// GetDevices returns an slice of all available Devices in the Device table from DynamoDB
// Returns a non-nil error if there's one during the execution and nil otherwise
func (db *DynamoDB) GetDevices() ([]types.Device, error) {
	out, err := db.dynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(db.DEVICES_TABLE_NAME),
	})

	if err != nil {
		err = fmt.Errorf("error getting information Devices table: %w", err)
		return nil, err
	}

	devices := []types.Device{}
	err = attributevalue.UnmarshalListOfMaps(out.Items, &devices)
	if err != nil {
		err = fmt.Errorf("error unmarshalling devices info: %w", err)
		return nil, err
	}

	return devices, nil
}

func (db *DynamoDB) InsertDevice(types.Device) error {
	return nil
}

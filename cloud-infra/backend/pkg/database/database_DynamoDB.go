package database

import (
	"backend/pkg/types"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	DynamoDBTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDB defines the struct used to implement Database interface using AWS DynamoDB
// It contains an DynamoDB client and the name of the tables to be used
type DynamoDB struct {
	dynamoDBClient    *dynamodb.Client
	DevicesTableName  string
	MessagesTableName string
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

	db.DevicesTableName = os.Getenv("DYNAMO_DB_DEVICES_TABLE_NAME")
	db.MessagesTableName = os.Getenv("DYNAMO_DB_MESSAGES_TABLE_NAME")

	db.dynamoDBClient = dynamodb.NewFromConfig(cfg)
}

// GetDevices returns an slice of all available Devices in the Device table from DynamoDB
// Returns a non-nil error if there's one during the execution and nil otherwise
func (db *DynamoDB) GetDevices() ([]types.Device, error) {
	out, err := db.dynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(db.DevicesTableName),
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

// InsertDevice receives a Device and inserts it in the Device table from DynamoDB
// Returns a non-nil error if there's one during the execution and nil otherwise
func (db *DynamoDB) InsertDevice(device types.Device) error {
	var err error
	if device.Model == "" {
		_, err = db.dynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(db.DevicesTableName),
			Item: map[string]DynamoDBTypes.AttributeValue{
				"DeviceUUID": &DynamoDBTypes.AttributeValueMemberS{Value: device.DeviceUUID},
				"Name":       &DynamoDBTypes.AttributeValueMemberS{Value: device.Name},
				"IP":         &DynamoDBTypes.AttributeValueMemberS{Value: device.IP},
			},
		})
	} else {
		_, err = db.dynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(db.DevicesTableName),
			Item: map[string]DynamoDBTypes.AttributeValue{
				"DeviceUUID": &DynamoDBTypes.AttributeValueMemberS{Value: device.DeviceUUID},
				"Name":       &DynamoDBTypes.AttributeValueMemberS{Value: device.Name},
				"IP":         &DynamoDBTypes.AttributeValueMemberS{Value: device.IP},
				"Model":      &DynamoDBTypes.AttributeValueMemberS{Value: device.Model},
			},
		})
	}
	if err != nil {
		err = fmt.Errorf("error while inserting: %w", err)
	}
	return err

}

// DeviceExistWithNameAndIP receives a device name and device ip and checks if there is any
// device that already have one of those 2 attributes matching exactly. Returns true is so and false otherwise
// Returns a non-nil error if there's one during the execution and nil otherwise
func (db *DynamoDB) DeviceExistWithNameAndIP(name string, ip string) (bool, error) {
	expr, err := expression.NewBuilder().WithFilter(
		expression.Or(
			expression.Equal(expression.Name("IP"), expression.Value(ip)),
			expression.Equal(expression.Name("Name"), expression.Value(name)),
		),
	).Build()
	if err != nil {
		err = fmt.Errorf("error while building the expression: %w", err)
		return true, err
	}

	out, err := db.dynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 aws.String(db.DevicesTableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		err = fmt.Errorf("error while scanning the DB: %w", err)
		return true, err
	}

	if len(out.Items) != 0 {
		return true, nil
	}
	return false, nil
}

func (db *DynamoDB) DeviceFromName(name string) (string, error) {
	expr, err := expression.NewBuilder().WithFilter(
		expression.Equal(expression.Name("Name"), expression.Value(name)),
	).Build()
	if err != nil {
		err = fmt.Errorf("error while building the expression: %w", err)
		return "", err
	}

	out, err := db.dynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 aws.String(db.DevicesTableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		err = fmt.Errorf("error while scanning the DB: %w", err)
		return "", err
	}

	if len(out.Items) == 0 {
		return "", nil
	}

	device := types.Device{}
	err = attributevalue.UnmarshalMap(out.Items[0], &device)
	if err != nil {
		err = fmt.Errorf("error unmarshalling device info: %w", err)
		return "", err
	}

	return device.IP, nil

}

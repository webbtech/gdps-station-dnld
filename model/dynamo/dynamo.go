package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pulpfree/gdps-fs-dwnld/config"
)

// Dynamo struct
type Dynamo struct {
	config *config.Dynamo
	db     *dynamodb.DynamoDB
}

// NewDB connection function
func NewDB(cfg *config.Dynamo) (*Dynamo, error) {

	var err error

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)

	return &Dynamo{
		config: cfg,
		db:     svc,
	}, err
}

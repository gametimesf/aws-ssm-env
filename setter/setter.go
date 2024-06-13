package setter

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type (
	Setter interface {
		Set(key string, value interface{})
	}

	service struct {
		client *ssm.SSM
	}
)

// NewSetter creates a new service that can store key value pairs in the parameter store
// ssmRegion: the region where the parameter store is located
func NewSetter(ssmRegion string) *service {
	awsConfig := aws.NewConfig()
	if ssmRegion != "" {
		awsConfig = awsConfig.WithRegion(ssmRegion)
	}

	session := session.Must(session.NewSession())
	client := ssm.New(session, awsConfig)

	return &service{
		client: client,
	}
}

// PutParameter stores/update the key value pair in the parameter store
// key: the name of the parameter
// value: the value of the parameter
// valueType: the type of the parameter. accepted types (String, StringList, SecureString)
// allowOverwrite: if true, the parameter will be overwritten if it already exists
func (s *service) PutParameter(ctx context.Context, key string, value string, valueType string, allowOverwrite bool) error {
	if valueType == "" || (valueType != "String" && valueType != "StringList" && valueType != "SecureString") {
		return fmt.Errorf("invalid value_type, use one of the expected values")
	}

	// Create/Update the parameter
	input := &ssm.PutParameterInput{
		Name:      aws.String(key),
		Value:     aws.String(value),
		Type:      aws.String(valueType),
		Overwrite: aws.Bool(allowOverwrite),
	}

	if _, err := s.client.PutParameterWithContext(ctx, input); err != nil {
		return fmt.Errorf("failed to put parameter %s: %v", key, err)
	}

	return nil
}

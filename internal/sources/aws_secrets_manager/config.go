package aws_secrets_manager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
	Region     string
	SecretName string
	VersionID  string
}

func (conf *Config) Fetch() (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(conf.Region),
	})
	if err != nil {
		return "", err
	}

	svc := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(conf.SecretName),
	}

	if conf.VersionID != "" {
		input.VersionId = aws.String(conf.VersionID)
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	if result.SecretString != nil {
		return *result.SecretString, nil
	}

	return "", nil
}

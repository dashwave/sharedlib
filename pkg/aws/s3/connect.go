package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dashwave/sharedlib/pkg/vault"

	sharedAws "github.com/dashwave/sharedlib/pkg/aws"
)

const (
	US_VAULT    = "US-VAULT"
	INDIA_VAULT = "INDIA-VAULT"
)

func ConnectAws(v *vault.VaultClient, region, accountLocation string) (*session.Session, vault.VaultSecretMap, error) {
	secretPath := ""
	if accountLocation == US_VAULT {
		secretPath = sharedAws.US_VAULT_SECRET_PATH
	} else if accountLocation == INDIA_VAULT {
		secretPath = sharedAws.INDIA_VAULT_SECRET_PATH
	} else {
		return nil, nil, fmt.Errorf("invalid AWS account location provided : %s", accountLocation)
	}
	secrets, err := v.GetSecretMap(secretPath)
	if err != nil {
		return nil, nil, err
	}
	accessKeyID := secrets[sharedAws.AWS_ACCESS_KEY_ID].(string)
	secretAccessKey := secrets[sharedAws.AWS_SECRET_ACCESS_KEY].(string)
	myRegion := region
	newSess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(myRegion),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		})
	session := session.Must(newSess, err)
	return session, secrets, nil
}

func ConnectS3(region string) (*session.Session, *s3.S3, vault.VaultSecretMap) {
	vaultClient, err := vault.NewVaultClient()
	if err != nil {
		panic(err)
	}
	awsSession, secrets, err := ConnectAws(vaultClient, region, INDIA_VAULT)
	if err != nil {
		panic(err)
	}
	s3Session := s3.New(awsSession)
	return awsSession, s3Session, secrets
}

func GetAWSSession(vaultToken, region, accountLocation string) (*session.Session, error) {
	vc, err := vault.GetVaultClientByToken(vaultToken)
	if err != nil {
		panic(err)
	}

	secretPath := ""
	if accountLocation == US_VAULT {
		secretPath = sharedAws.US_VAULT_SECRET_PATH
	} else if accountLocation == INDIA_VAULT {
		secretPath = sharedAws.INDIA_VAULT_SECRET_PATH
	} else {
		return nil, fmt.Errorf("invalid AWS account location provided : %s", accountLocation)
	}
	secrets, err := vc.GetSecretMapByStore(secretPath, sharedAws.AWS_CREDENTIALS_STORE)
	if err != nil {
		return nil, err
	}
	accessKeyID := secrets[sharedAws.AWS_ACCESS_KEY_ID].(string)
	secretAccessKey := secrets[sharedAws.AWS_SECRET_ACCESS_KEY].(string)
	newSess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		})
	session := session.Must(newSess, err)
	return session, nil
}

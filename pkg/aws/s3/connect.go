package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dashwave/sharedlib/pkg/vault"

	sharedAws "github.com/dashwave/sharedlib/pkg/aws"
)

func connectAws(v *vault.VaultClient, region string) (*session.Session, vault.VaultSecretMap, error) {
	secrets, err := v.GetSecretMap(sharedAws.VAULT_SECRET_PATH)
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

func ConnectS3(region string) (*s3.S3, vault.VaultSecretMap) {
	vaultClient, err := vault.NewVaultClient()
	if err != nil {
		panic(err)
	}
	awsSession, secrets, err := connectAws(vaultClient, region)
	if err != nil {
		panic(err)
	}
	s3Session := s3.New(awsSession)
	return s3Session, secrets
}

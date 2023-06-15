package s3

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dashwave/sharedlib/pkg/vault"

	sharedAws "github.com/dashwave/sharedlib/pkg/aws"
)

// CreateBucket creates a new S3 bucket with the specified bucketname using the provided S3 session.
// If a bucket with the provided name already exist on our account, it would return stating the log for
// the same. Throws an error if bucket exist on an account not owned by request user, since bucket names are
// of global namespace. If a new bucket is created, this function also enables ACL and versioning on the new bucket.
func CreateBucket(sess *s3.S3, secrets vault.VaultSecretMap, config *CreateBucketConfiguration) error {
	region := secrets[sharedAws.AWS_REGION_KEY].(string)
	createBucketRequest := &s3.CreateBucketInput{
		Bucket: aws.String(config.Name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(region),
		},
	}
	if _, err := sess.CreateBucket(createBucketRequest); err != nil {
		switch strings.Split(err.Error(), ":")[0] {
		case s3.ErrCodeBucketAlreadyOwnedByYou:
			fmt.Printf("Bucket with the name %s already exists in our account, using the existing bucket\n", config.Name)
			return nil
		default:
			return err
		}
	}
	fmt.Printf("Successfully created new bucket with name %v\n", config.Name)

	if config.EnableACL {
		if err := enableBucketACL(sess, config.Name); err != nil {
			return err
		}
	}
	if config.EnableVersionsing {
		if err := enableBucketVersioning(sess, config.Name); err != nil {
			return err
		}
	}
	if config.EnableTransferAcceleration {
		if err := enableBucketAccelerateTransfer(sess, config.Name); err != nil {
			return err
		}
	}
	return nil
}

// enableBucketACL enables ACL on the bucket specified. AWS turns off ACL by default and blocks all possiblity
// of public read through ACL. We are disabling the blocking, and enabling the ACL. The bucket access is only
// given to the user, but bucket access can be controlled later by setting ACL to change permissions.
func enableBucketACL(sess *s3.S3, bucketName string) error {
	reqOwnership := &s3.PutBucketOwnershipControlsInput{
		Bucket: aws.String(bucketName),
		OwnershipControls: &s3.OwnershipControls{
			Rules: []*s3.OwnershipControlsRule{
				{
					ObjectOwnership: aws.String("ObjectWriter"),
				},
			},
		},
	}
	_, err := sess.PutBucketOwnershipControls(reqOwnership)
	if err != nil {
		return err
	}
	reqAccess := &s3.PutPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
		PublicAccessBlockConfiguration: &s3.PublicAccessBlockConfiguration{
			BlockPublicPolicy: aws.Bool(false),
		},
	}
	_, err = sess.PutPublicAccessBlock(reqAccess)
	if err != nil {
		return err
	}
	fmt.Println("Successfully enabled ACL on bucket")
	return nil
}

// enableBucketVersioning enables versioning sysytem on the bucket with the given bucketname. If multiple
// objects are uploaded to this bucket with the same key value, all the versions of that object are stored,
// with the most recent one set as the default.
func enableBucketVersioning(sess *s3.S3, bucketName string) error {
	versioningReq := s3.PutBucketVersioningInput{
		Bucket: aws.String(bucketName),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String("Enabled"),
		},
	}
	if _, err := sess.PutBucketVersioning(&versioningReq); err != nil {
		return err
	}
	fmt.Println("Successfully enabled versioning in bucket")
	return nil
}

// enableBucketAccelerateTransfer enables aws feature of Transfer Acclereration to upload file
// to a nearest edge location and then route it to the final destination via an optimised path.
func enableBucketAccelerateTransfer(sess *s3.S3, bucketName string) error {
	req := &s3.PutBucketAccelerateConfigurationInput{
		Bucket: aws.String(bucketName),
		AccelerateConfiguration: &s3.AccelerateConfiguration{
			Status: aws.String(s3.BucketAccelerateStatusEnabled),
		},
	}
	_, err := sess.PutBucketAccelerateConfiguration(req)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully enabled accelerated transfer for bucket: %v\n", bucketName)
	return nil
}

// UploadObjectToBucket uploads the provided object to S3 bucket. It either completely uploads the object to the bucket
// and returns successfully or throws an error without any upload.
func UploadObjectToBucket(sess *s3.S3, object *S3Object) error {
	objectReq := &s3.PutObjectInput{
		Bucket: object.Bucket,
		Key:    object.Key,
		Body:   bytes.NewReader(object.Body),
		ACL:    aws.String(object.ACL),
	}

	if err := objectReq.Validate(); err != nil {
		return err
	}

	_, err := sess.PutObject(objectReq)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded object with key %v to the bucket %v.\n", *object.Key, *object.Bucket)
	return nil
}

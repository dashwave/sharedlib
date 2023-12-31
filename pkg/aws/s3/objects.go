package s3

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadObjectToBucket uploads the provided object to S3 bucket. It either completely uploads the object to the bucket
// and returns successfully or throws an error without any upload.
func UploadObjectToBucket(s3Session *s3.S3, object *S3Object) error {
	objectReq := &s3.PutObjectInput{
		Bucket: object.Bucket,
		Key:    object.Key,
		Body:   bytes.NewReader(object.Body),
		ACL:    aws.String(object.ACL),
	}

	if err := objectReq.Validate(); err != nil {
		return err
	}

	_, err := s3Session.PutObject(objectReq)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded object with key %v to the bucket %v.\n", *object.Key, *object.Bucket)
	return nil
}

// UploadObjectMultipart uploads the object data from the given source object to the bucket.
// This is achieved by dividing the data into multiple parts and uploading them over
// concurrent streams which is by default set to 5.
// Set the desired location of source object/file along with key
func UploadObjectMultipart(awsSess *session.Session, r *UploadMultipartObjectRequest) error {
	file, err := os.Open(r.Source)
	if err != nil {
		fmt.Println("Error opening local file:", err)
		return err
	}
	defer file.Close()

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(r.ObjectName),
		Body:   file,
	}

	uploader := s3manager.NewUploader(awsSess, func(d *s3manager.Uploader) {
		d.PartSize = 200 * 1024 * 1024 // 200MB per part
	})

	_, err = uploader.Upload(upParams)
	if err != nil {
		return err
	}

	return nil
}

// UploadObjectMultipartWithContext uploads the object data from the given source object to the bucket.
// This is achieved by dividing the data into multiple parts and uploading them over
// concurrent streams which is by default set to 5.
// Set the desired location of source object/file along with key
// Takes in a context to stop the request when context is expired
func UploadObjectMultipartWithContext(ctx context.Context, awsSess *session.Session, r *UploadMultipartObjectRequest) error {
	file, err := os.Open(r.Source)
	if err != nil {
		fmt.Println("Error opening local file:", err)
		return err
	}
	defer file.Close()

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(r.ObjectName),
		Body:   file,
	}

	uploader := s3manager.NewUploader(awsSess, func(d *s3manager.Uploader) {
		d.PartSize = 200 * 1024 * 1024 // 200MB per part
	})

	_, err = uploader.UploadWithContext(ctx, upParams)
	if err != nil {
		return err
	}

	return nil
}

// GetObject downloads the object data for the given object key from the bucket. To get an object with a
// specific version id, set VersioningEnabled to true and provide the version id.
func GetObject(s3Session *s3.S3, r *GetObjectRequest) (*GetObjectResponse, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(r.ObjectName),
	}
	if r.VersioningEnabled {
		getObjectInput.VersionId = aws.String(r.VersionId)
	}
	resp, err := s3Session.GetObject(getObjectInput)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	object := &GetObjectResponse{
		Body: data,
	}
	return object, nil
}

// GetObjectMultiipart downloads the object data for the given object key from the bucket.
// This is achieved by dividing the data into multiple parts and downloading them over
// concurrent steams which is by default set to 5.
// Set the desired location of downloaded data with destination
func GetObjectMultipart(awsSess *session.Session, r *GetMultiPartObjectRequest) error {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(r.ObjectName),
	}
	if r.VersioningEnabled {
		getObjectInput.VersionId = aws.String(r.VersionId)
	}

	file, err := os.Create(r.Destination)
	if err != nil {
		return err
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(awsSess, func(d *s3manager.Downloader) {
		d.PartSize = 200 * 1024 * 1024 // 200MB per part
	})

	_, err = downloader.Download(file, getObjectInput)
	if err != nil {
		return err
	}

	return nil
}

// DoesObjectExists checks if a particular object exist in the specified bucket
// and returns corresponding boolean value
func DoesObjectExists(s3Session *s3.S3, r *ObjectExistsReq) (bool, error) {
	_, err := s3Session.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(r.ObjectName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}

	return true, nil
}

// DoesObjectsWithPrefix return the list of objects that exist with the
// given prefix
func ListObjectsWithPrefix(s3Session *s3.S3, r *ListObjectsReq) ([]*s3.Object, error) {
	res, err := s3Session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:  aws.String(r.BucketName),
		Prefix:  aws.String(r.Prefix),
		MaxKeys: aws.Int64(r.MaxKeys),
	})
	if err != nil {
		return nil, err
	}

	return res.Contents, nil
}

// GetObjectPresignedURL generates the public URL to download the object data for the given object key from the private bucket.
// To get an object with aspecific version id, set VersioningEnabled to true and provide the version id.
// Returns the public URL, which is valid for specific Duration given in request
func GetObjectPresignedURL(s3Session *s3.S3, r *GetObjectRequest) (string, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(r.ObjectName),
	}
	if r.VersioningEnabled {
		getObjectInput.VersionId = aws.String(r.VersionId)
	}

	req, _ := s3Session.GetObjectRequest(getObjectInput)

	urlStr, err := req.Presign(r.Duration)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}

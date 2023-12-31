package s3

import "time"

type CreateBucketConfiguration struct {
	Name                       string
	EnableVersionsing          bool
	EnableACL                  bool
	EnableTransferAcceleration bool
}

type S3Object struct {
	Bucket *string
	Key    *string
	Body   []byte
	ACL    string
}

type GetObjectRequest struct {
	BucketName        string
	ObjectName        string
	VersioningEnabled bool
	VersionId         string
	Duration          time.Duration
}

type GetObjectResponse struct {
	Body []byte
}

type GetMultiPartObjectRequest struct {
	BucketName        string
	ObjectName        string
	VersioningEnabled bool
	VersionId         string
	Destination       string
}

type ObjectExistsReq struct {
	BucketName string
	ObjectName string
}

type ListObjectsReq struct {
	BucketName string
	Prefix     string
	MaxKeys    int64
}

type UploadMultipartObjectRequest struct {
	BucketName string
	ObjectName string
	Source     string
}

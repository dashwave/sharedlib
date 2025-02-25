package storage

import "time"

type CreateBucketConfiguration struct {
	Name                string
	Location            string
	EnableVersioning    bool
	EnableUniformAccess bool
}

type StorageObject struct {
	Bucket *string
	Name   *string
	Body   []byte
	ACL    string
}

type GetObjectRequest struct {
	BucketName        string
	ObjectName        string
	VersioningEnabled bool
	Generation        int64
	Duration          time.Duration
}

type GetObjectResponse struct {
	Body        []byte
	ContentType string
}

type GetMultiPartObjectRequest struct {
	BucketName        string
	ObjectName        string
	VersioningEnabled bool
	Generation        int64
	Destination       string
	Source            string
}

type ObjectExistsReq struct {
	BucketName string
	ObjectName string
}

type ListObjectsReq struct {
	BucketName string
	Prefix     string
	MaxResults int64
}

package s3

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
}

type GetObjectResponse struct {
	Body []byte
}

type GetMultiPartObjectRequest struct {
	BucketName        string
	ObjectName        string
	VersioningEnabled bool
	VersionId         string
	destination       string
}

type ObjectExistsReq struct {
	BucketName string
	ObjectName string
}

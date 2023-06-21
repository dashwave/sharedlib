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

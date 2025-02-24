package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

// CreateBucket creates a new GCS bucket with the specified configuration.
// If a bucket with the provided name already exists in our project, it returns stating the log for
// the same. If a new bucket is created, this function also enables uniform bucket-level access and versioning
// based on the configuration.
func CreateBucket(client *storage.Client, config *CreateBucketConfiguration) error {
	ctx := context.Background()
	bucket := client.Bucket(config.Name)

	// Check if bucket already exists
	_, err := bucket.Attrs(ctx)
	if err == nil {
		fmt.Printf("Bucket with the name %s already exists in our project, using the existing bucket\n", config.Name)
		return nil
	}
	if err != storage.ErrBucketNotExist {
		return err
	}

	// Create bucket with specified location
	attrs := &storage.BucketAttrs{
		Location: config.Location,
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: config.EnableUniformAccess,
		},
	}
	if err := bucket.Create(ctx, "", attrs); err != nil {
		return err
	}
	fmt.Printf("Successfully created new bucket with name %v\n", config.Name)

	if config.EnableVersioning {
		if err := enableBucketVersioning(client, config.Name); err != nil {
			return err
		}
	}

	return nil
}

// enableBucketVersioning enables versioning system on the bucket with the given bucketname.
// If multiple objects are uploaded to this bucket with the same name, all the versions of that object are stored,
// with the most recent one set as the default.
func enableBucketVersioning(client *storage.Client, bucketName string) error {
	ctx := context.Background()
	bucket := client.Bucket(bucketName)

	update := storage.BucketAttrsToUpdate{
		VersioningEnabled: true,
	}
	if _, err := bucket.Update(ctx, update); err != nil {
		return err
	}
	fmt.Println("Successfully enabled versioning in bucket")
	return nil
}

// DeleteBucket deletes the bucket with the provided bucketname
func DeleteBucket(client *storage.Client, bucketName string) error {
	ctx := context.Background()
	bucket := client.Bucket(bucketName)

	if err := bucket.Delete(ctx); err != nil {
		return err
	}
	return nil
}

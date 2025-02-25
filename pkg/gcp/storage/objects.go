package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// UploadObjectToBucket uploads the provided object to GCS bucket. It either completely uploads the object
// to the bucket and returns successfully or throws an error without any upload.
func UploadObjectToBucket(client *storage.Client, object *StorageObject) error {
	ctx := context.Background()
	bucket := client.Bucket(*object.Bucket)
	obj := bucket.Object(*object.Name)

	writer := obj.NewWriter(ctx)
	if object.ACL != "" {
		writer.PredefinedACL = object.ACL
	}

	if _, err := writer.Write(object.Body); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded object with name %v to the bucket %v.\n", *object.Name, *object.Bucket)
	return nil
}

// UploadObjectMultipart uploads the object data from the given source object to the bucket.
// GCS automatically handles chunking and parallel uploads for large files.
func UploadObjectMultipart(client *storage.Client, r *GetMultiPartObjectRequest) error {
	ctx := context.Background()
	file, err := os.Open(r.Source)
	if err != nil {
		fmt.Println("Error opening local file:", err)
		return err
	}
	defer file.Close()

	bucket := client.Bucket(r.BucketName)
	obj := bucket.Object(r.ObjectName)

	writer := obj.NewWriter(ctx)
	if _, err := io.Copy(writer, file); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

// GetObject downloads the object data for the given object name from the bucket.
// To get an object with a specific generation, set VersioningEnabled to true and provide the generation number.
func GetObject(client *storage.Client, r *GetObjectRequest) (*GetObjectResponse, error) {
	ctx := context.Background()
	bucket := client.Bucket(r.BucketName)
	obj := bucket.Object(r.ObjectName)

	if r.VersioningEnabled && r.Generation > 0 {
		obj = obj.Generation(r.Generation)
	}

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	return &GetObjectResponse{Body: data, ContentType: attrs.ContentType}, nil
}

// GetObjectMultipart downloads the object data for the given object name from the bucket.
// GCS automatically handles chunking and parallel downloads for large files.
func GetObjectMultipart(client *storage.Client, r *GetMultiPartObjectRequest) error {
	ctx := context.Background()
	bucket := client.Bucket(r.BucketName)
	obj := bucket.Object(r.ObjectName)

	if r.VersioningEnabled && r.Generation > 0 {
		obj = obj.Generation(r.Generation)
	}

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	file, err := os.Create(r.Destination)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return err
	}

	return nil
}

// DoesObjectExists checks if a particular object exists in the specified bucket
// and returns corresponding boolean value
func DoesObjectExists(client *storage.Client, r *ObjectExistsReq) (bool, error) {
	ctx := context.Background()
	bucket := client.Bucket(r.BucketName)
	obj := bucket.Object(r.ObjectName)

	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// ListObjectsWithPrefix returns the list of objects that exist with the
// given prefix
func ListObjectsWithPrefix(client *storage.Client, r *ListObjectsReq) ([]*storage.ObjectAttrs, error) {
	ctx := context.Background()
	bucket := client.Bucket(r.BucketName)

	var objects []*storage.ObjectAttrs
	it := bucket.Objects(ctx, &storage.Query{
		Prefix: r.Prefix,
	})

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		objects = append(objects, attrs)
	}

	return objects, nil
}

// GetObjectSignedURL generates a signed URL to download the object data for the given object name from the bucket.
// To get an object with a specific generation, set VersioningEnabled to true and provide the generation number.
// Returns the signed URL, which is valid for specific Duration given in request
func GetObjectSignedURL(client *storage.Client, r *GetObjectRequest) (string, error) {
	bucket := client.Bucket(r.BucketName)

	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(r.Duration),
	}

	return bucket.SignedURL(r.ObjectName, opts)
}

// GetUploadSignedURL generates a signed URL that can be used to upload an object to the bucket.
// The URL will be valid for the specified duration.
func GetUploadSignedURL(client *storage.Client, r *GetObjectRequest) (string, error) {
	bucket := client.Bucket(r.BucketName)

	opts := &storage.SignedURLOptions{
		Method:  "PUT",
		Expires: time.Now().Add(r.Duration),
	}

	return bucket.SignedURL(r.ObjectName, opts)
}

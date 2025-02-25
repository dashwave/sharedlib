package storage

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
)

const (
	testBucketName = "test-bucket-1"
	testObjectName = "test-object.txt"
	testContent    = "Hello, GCP Storage!"
)

func setupTestClient(t *testing.T) *storage.Client {
	// Ensure GOOGLE_APPLICATION_CREDENTIALS environment variable is set
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		t.Fatalf("Failed to create storage client: %v", err)
	}
	return client
}

func TestCreateAndDeleteBucket(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	config := &CreateBucketConfiguration{
		Name:                testBucketName,
		Location:            "US-CENTRAL1",
		EnableVersioning:    true,
		EnableUniformAccess: true,
	}

	// Create bucket
	err := CreateBucket(client, config)
	assert.NoError(t, err)

	// Delete bucket
	// err = DeleteBucket(client, testBucketName)
	// assert.NoError(t, err)
}

func TestObjectOperations(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	// Create test bucket
	// config := &CreateBucketConfiguration{
	// 	Name:                testBucketName,
	// 	Location:            "US-CENTRAL1",
	// 	EnableVersioning:    true,
	// 	EnableUniformAccess: true,
	// }
	// err := CreateBucket(client, config)
	// assert.NoError(t, err)
	// defer DeleteBucket(client, testBucketName)

	// Test upload - Remove ACL since uniform access is enabled
	bucketName := testBucketName
	objectName := testObjectName
	object := &StorageObject{
		Bucket: &bucketName,
		Name:   &objectName,
		Body:   []byte(testContent),
	}
	err := UploadObjectToBucket(client, object)
	assert.NoError(t, err)

	// Add error checking and potential retry/delay
	if err == nil {
		// Test object exists
		exists, err := DoesObjectExists(client, &ObjectExistsReq{
			BucketName: testBucketName,
			ObjectName: testObjectName,
		})
		assert.NoError(t, err)
		assert.True(t, exists)

		// Test download
		downloadReq := &GetObjectRequest{
			BucketName: testBucketName,
			ObjectName: testObjectName,
		}
		resp, err := GetObject(client, downloadReq)
		assert.NoError(t, err)
		if err == nil {
			assert.Equal(t, testContent, string(resp.Body))
		}

		// Test list objects
		listReq := &ListObjectsReq{
			BucketName: testBucketName,
			Prefix:     "test",
		}
		objects, err := ListObjectsWithPrefix(client, listReq)
		assert.NoError(t, err)
		assert.Len(t, objects, 1)
		assert.Equal(t, testObjectName, objects[0].Name)

		// Test signed URL
		urlReq := &GetObjectRequest{
			BucketName: testBucketName,
			ObjectName: testObjectName,
			Duration:   time.Hour,
		}
		signedURL, err := GetObjectSignedURL(client, urlReq)

		fmt.Println("signedURL", signedURL)
		assert.NoError(t, err)
		assert.Contains(t, signedURL, testBucketName)
		assert.Contains(t, signedURL, testObjectName)
	}
}

func TestMultipartOperations(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	// Create test file
	tempFile := "test_upload.txt"
	err := os.WriteFile(tempFile, []byte(testContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(tempFile)

	// Create test bucket
	// config := &CreateBucketConfiguration{
	// 	Name:                testBucketName,
	// 	Location:            "US-CENTRAL1",
	// 	EnableVersioning:    true,
	// 	EnableUniformAccess: true,
	// }
	// err = CreateBucket(client, config)
	// assert.NoError(t, err)
	// defer DeleteBucket(client, testBucketName)

	// Test multipart upload
	uploadReq := &GetMultiPartObjectRequest{
		BucketName: testBucketName,
		ObjectName: "multipart-test.txt",
		Source:     tempFile,
	}
	err = UploadObjectMultipart(client, uploadReq)
	assert.NoError(t, err)

	// Test multipart download
	downloadFile := "test_download.txt"
	downloadReq := &GetMultiPartObjectRequest{
		BucketName:  testBucketName,
		ObjectName:  "multipart-test.txt",
		Destination: downloadFile,
	}
	err = GetObjectMultipart(client, downloadReq)
	assert.NoError(t, err)
	defer os.Remove(downloadFile)

	// Verify downloaded content
	content, err := os.ReadFile(downloadFile)
	assert.NoError(t, err)
	assert.Equal(t, testContent, string(content))
}

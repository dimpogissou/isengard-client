package connectors

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
)

// Util function creating S3 bucket for integration test,
// does not require deletion post-test since Localstack API calls persistence is disabled on this project
func createS3Bucket(client *s3.S3, bucket string) {

	_, err := client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		logger.Error("CannotCreateTestS3Bucket", fmt.Sprintf("%v", err))
		panic(fmt.Sprintf("Unable to create test bucket: %v", err))
	}

	// Wait until bucket is created before finishing
	logger.Info(fmt.Sprintf("Waiting for bucket %q to be created...\n", bucket))

	err = client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
}

func listFilesFromS3Bucket(client *s3.S3, bucket string, prefix string) []*s3.Object {

	listObjectsInput := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	resp, err := client.ListObjects(listObjectsInput)
	if err != nil {
		panic(err.Error())
	}

	return resp.Contents
}

func readFileFromS3Bucket(downloader *s3manager.Downloader, bucket string, key string) string {

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		panic(err.Error())
	}

	return string(buf.Bytes())
}

func deleteFilesAndBucket(client *s3.S3, bucket string, prefix string) {

	// Delete bucket items: not using batch delete due to an issue with Localstack
	// Using Bucket list and iterating deletes as a workaround
	// See https://github.com/localstack/localstack/issues/1452#issuecomment-711422073

	// List files from bucket
	s3Files := listFilesFromS3Bucket(client, bucket, prefix)

	// Delete files from bucket
	for _, f := range s3Files {
		_, err := client.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: f.Key})
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Successfully deleted files from bucket %s", bucket)

	// Delete bucket itself

	_, err := client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		panic(err.Error())
	}

	// Wait until bucket is deleted before finishing
	fmt.Printf("Waiting for bucket %q to be deleted...\n", bucket)

	err = client.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	fmt.Printf("Successfully deleted bucket %s", bucket)

}

func TestWriteToS3(t *testing.T) {

	localstackEndpoint := "http://localstack"
	localstackS3Port := "4572"
	testEndpoint := fmt.Sprintf("%s:%s", localstackEndpoint, localstackS3Port)
	testRegion := os.Getenv("AWS_DEFAULT_REGION")

	const testBucket = "local-test-bucket"
	const testKeyPrefix = "some/test/key"
	const testString = "Log line"
	const nFiles = 5

	// Create test S3 connector
	cfg := S3ConnectorConfig{
		Name:      "testS3Connector",
		Type:      "s3",
		Endpoint:  testEndpoint,
		KeyPrefix: testKeyPrefix,
		Region:    testRegion,
		Bucket:    testBucket,
		Levels:    []string{"INFO", "DEBUG", "WARN", "ERROR"},
	}
	session, client := SetupS3Client(cfg)
	connector := S3Connector{cfg: cfg, session: session, client: client}

	// Create test S3 bucket and defer its deletion
	createS3Bucket(connector.client, testBucket)
	defer deleteFilesAndBucket(connector.client, testBucket, testKeyPrefix)

	// No-op at the moment
	connector.Open()
	defer connector.Close()

	// Send lines to test s3 bucket
	line := tail.Line{Text: testString}
	for i := 0; i < nFiles; i++ {
		if connector.Send(&line) != true {
			t.Errorf("S3 connector Send function failed")
		}
	}

	// List files in bucket
	s3Files := listFilesFromS3Bucket(connector.client, testBucket, testKeyPrefix)
	if len(s3Files) != nFiles {
		t.Errorf("Number of files in S3 not matching expectation: got %d, want %d", len(s3Files), nFiles)
	}

	// Download files from bucket and run assertions
	downloader := s3manager.NewDownloader(connector.session)
	for _, f := range s3Files {
		fileContent := readFileFromS3Bucket(downloader, testBucket, *f.Key)
		if fileContent != testString {
			t.Errorf("S3 file content not matching expectation, got %s, want %s", fileContent, testString)
		}
	}

}

package connectors

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dimpogissou/isengard-server/config"
	"github.com/hpcloud/tail"
)

func createS3Bucket(client *s3.S3, bucket string) {

	_, err := client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Error(err)
		panic("Unable to create test bucket, failing")
	}

	// Wait until bucket is created before finishing
	log.Info("Waiting for bucket %q to be created...\n", bucket)

	err = client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
}

func TestSend(t *testing.T) {

	const testEndpoint = "http://localstack:4572"
	const testRegion = "eu-west-1a"
	const testBucket = "local-test-bucket"
	const testKeyPrefix = "some/test/key/"
	cfg := config.Connector{
		Name:      "testS3Connector",
		Type:      "s3",
		Endpoint:  testEndpoint,
		KeyPrefix: testKeyPrefix,
		Region:    testRegion,
		Bucket:    testBucket,
		Levels:    []string{"INFO", "DEBUG", "WARN", "ERROR"},
	}
	connector := S3Connector{cfg: cfg, client: SetupS3Client(cfg)}
	createS3Bucket(connector.client, testBucket)

	connector.Open()
	defer connector.Close()

	// Send line to test s3 bucket
	line := tail.Line{Text: "Log line"}
	if connector.Send(&line) != true {
		t.Errorf("S3 connector Send function failed")
	}

}

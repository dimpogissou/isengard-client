package connectors

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
	uuid "github.com/nu7hatch/gouuid"
)

type S3Connector struct {
	session *session.Session
	client  *s3.S3
	cfg     S3ConnectorConfig
}

// Sets up S3 client
func SetupS3Client(cfg S3ConnectorConfig) (*session.Session, *s3.S3) {
	sessionPtr := session.Must(session.NewSession(&aws.Config{
		S3ForcePathStyle:              aws.Bool(true),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Region:                        aws.String(cfg.Region),
		Endpoint:                      aws.String(cfg.Endpoint),
	}))
	client := s3.New(sessionPtr, &aws.Config{})
	logger.Info(fmt.Sprintf("Created S3 client --> %v", &client))
	return sessionPtr, client
}

func (c S3Connector) Close() error {
	logger.Info("Closed S3 connector (no-op) ...")
	return nil
}

// Parses a log line into a string map using the regex built from config
func parseLine(l *tail.Line, re *regexp.Regexp) (map[string]string, error) {
	match := re.FindStringSubmatch(l.Text)
	if match == nil {
		return make(map[string]string), errors.New("No match found in line, returning empty map")
	} else {
		paramsMap := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}
		return paramsMap, nil
	}
}

// Puts a tailed line into the specified bucket
func (c S3Connector) s3PutObject(bucket string, fileKey string, line *tail.Line) (*s3.PutObjectOutput, error) {

	p := s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		ACL:    aws.String("public-read"),
		Body:   strings.NewReader(line.Text),
	}

	r, err := c.client.PutObject(&p)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c S3Connector) Send(line *tail.Line) error {
	t := time.Now()
	uuid, err := uuid.NewV4()
	if err != nil {
		logger.Error("CreateUuidError", err.Error())
		return err
	}
	fileName := fmt.Sprintf("%s/%d-%02d-%02dT%02d-%02d-%02d-%v",
		c.cfg.KeyPrefix,
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), uuid)
	_, err = c.s3PutObject(c.cfg.Bucket, fileName, line)
	logger.Info(fmt.Sprintf("Sending file '%s' to S3 bucket '%s'", fileName, c.cfg.Bucket))
	if err != nil {
		logger.Error("S3PutObjectError", err.Error())
		return err
	}
	return nil
}

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
	"github.com/dimpogissou/isengard-server/config"
	"github.com/hpcloud/tail"
	uuid "github.com/nu7hatch/gouuid"
)

type S3Connector struct {
	client *s3.S3
	cfg    config.Connector
}

func SetupS3Client(cfg config.Connector) *s3.S3 {
	sessionPtr := session.Must(session.NewSession(&aws.Config{
		S3ForcePathStyle:              aws.Bool(true),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Region:                        aws.String(cfg.Region),
		Endpoint:                      aws.String(cfg.Endpoint),
	}))
	client := s3.New(sessionPtr, &aws.Config{})
	log.Info("Created S3 client -->", client)
	return client
}

func (conn S3Connector) Open() {
	log.Info("Starting S3 connector ...")
}

func (c S3Connector) Close() {
	log.Info("Closing S3 connector ...")
}

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

func (c S3Connector) s3PutObject(fileKey string, line *tail.Line) (*s3.PutObjectOutput, error) {

	p := s3.PutObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
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

// TODO -> return (struct, err)
func (c S3Connector) Send(line *tail.Line) bool {
	t := time.Now()
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Error("Failed generating UUID for S3 file, exiting Send function")
		return false
	}
	fileName := fmt.Sprintf("%s/%d-%02d-%02dT%02d-%02d-%02d-%v",
		c.cfg.KeyPrefix,
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), uuid)
	_, err = c.s3PutObject(fileName, line)
	log.Info(fmt.Sprintf("Sending file '%s' to S3 bucket '%s'", fileName, c.cfg.Bucket))
	if err != nil {
		return false
	}
	return true
}

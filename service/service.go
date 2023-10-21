package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/oleksiivelychko/go-aws-storage/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	bucketPrefix             = "s3://"
	defaultProtectedRegion   = "us-east-1"
	MinutesToExpireSignedURL = 15
)

type IService interface {
	ListBuckets() (string, error)
	CreateBucket(name string, isPublic bool) (string, error)
	PutObjectsAsync(bucket string, filenames []string, outCh chan string, errCh chan error)
	PutObjects(bucket string, filePaths []string) (output []string, errArr []error)
	ListObjects(bucket string) (string, error)
	GetObject(bucket, key, path string) error
	DeleteObject(bucket, key string) (string, error)
	DeleteBucket(bucket string) (string, error)
	AssignURL(bucket, key string) (string, error)
}

type service struct {
	client *s3.S3
}

func New(config *config.AWS) (IService, error) {
	awsConfig := &aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AwsAccessKeyId, config.AwsSecretAccessKey, ""),
	}

	if config.Endpoint != "" {
		awsConfig.Endpoint = aws.String(config.Endpoint)
	}

	if config.S3ForcePathStyle {
		awsConfig.S3ForcePathStyle = aws.Bool(config.S3ForcePathStyle)
	}

	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	return &service{client: s3.New(awsSession)}, nil
}

func (service *service) ListBuckets() (string, error) {
	output, err := service.client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func (service *service) CreateBucket(name string, isPublic bool) (string, error) {
	match, _ := regexp.MatchString("^[a-z0-9\\-]{3,50}$", name)
	if !match {
		return "", fmt.Errorf("bucket name \"%s\" is not applicable", name)
	}

	var bucketConfig *s3.CreateBucketConfiguration
	if *service.client.Config.Region != "us-east-1" {
		bucketConfig.LocationConstraint = service.client.Config.Region
	}

	input := &s3.CreateBucketInput{
		Bucket:                    aws.String(name),
		CreateBucketConfiguration: bucketConfig,
	}

	if isPublic {
		input.ACL = aws.String(s3.BucketCannedACLPublicRead)
	} else {
		input.ACL = aws.String(s3.BucketCannedACLPrivate)
	}

	output, err := service.client.CreateBucket(input)

	if err != nil {
		var awsErr awserr.Error

		if errors.As(err, &awsErr) {
			switch awsErr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				return "", fmt.Errorf("%s: %s", s3.ErrCodeBucketAlreadyExists, awsErr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				return "", fmt.Errorf("%s: %s", s3.ErrCodeBucketAlreadyOwnedByYou, awsErr.Error())
			default:
				return "", awsErr
			}
		}
	}

	return output.String(), err
}

func (service *service) PutObjectsAsync(bucket string, filenames []string, outCh chan string, errCh chan error) {
	var wg sync.WaitGroup

	for _, key := range filenames {
		wg.Add(1)

		go func(filename string) {
			defer wg.Done()

			output, err := service.putObjects(bucket, filename)
			if err != nil {
				errCh <- err
			} else {
				outCh <- fmt.Sprintf("file %s was uploaded and available by ETag %s\n", filename, *output.ETag)
			}
		}(key)
	}

	wg.Wait()

	close(outCh)
	close(errCh)
}

func (service *service) PutObjects(bucket string, filePaths []string) (output []string, errArr []error) {
	for _, filePath := range filePaths {
		result, err := service.putObjects(bucket, filePath)
		if err != nil {
			errArr = append(errArr, err)
		} else {
			output = append(output, fmt.Sprintf("file %s was uploaded\n%s\n", filePath, result.String()))
		}
	}

	return
}

func (service *service) ListObjects(bucket string) (string, error) {
	output, err := service.client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(10000),
	})

	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func (service *service) putObjects(bucket, filePath string) (*s3.PutObjectOutput, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := fileInfo.Size()
	buf := make([]byte, size)

	_, readErr := file.Read(buf)
	if readErr != nil {
		return nil, readErr
	}

	return service.client.PutObject(&s3.PutObjectInput{
		ACL:           aws.String(s3.BucketCannedACLPublicRead),
		Body:          aws.ReadSeekCloser(bytes.NewReader(buf)),
		Bucket:        aws.String(bucket),
		ContentType:   aws.String(http.DetectContentType(buf)),
		ContentLength: aws.Int64(size),
		Key:           aws.String(strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]),
	})
}

func (service *service) GetObject(bucket, key, path string) error {
	output, err := service.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	body, err := io.ReadAll(output.Body)
	if err != nil {
		return err
	}

	if path != "" {
		absPath, absPathErr := filepath.Abs(path)
		if absPathErr != nil {
			return absPathErr
		}

		if _, statErr := os.Stat(absPath); os.IsNotExist(statErr) {
			if err = os.Mkdir(absPath, os.ModePerm); err != nil {
				return err
			}
		}

		path = filepath.Join(absPath, key)
	} else {
		path = key
	}

	return os.WriteFile(path, body, 0644)
}

func (service *service) DeleteObject(bucket, key string) (string, error) {
	_, err := service.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("object %s was deleted from %s%s", key, bucketPrefix, bucket), nil
}

func (service *service) DeleteBucket(bucket string) (string, error) {
	_, err := service.client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("bucket %s%s was deleted", bucketPrefix, bucket), nil
}

func (service *service) AssignURL(bucket, key string) (string, error) {
	req, _ := service.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return req.Presign(MinutesToExpireSignedURL * time.Minute)
}

package awsutil

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gangjun06/d4dj-crawler/conf"
)

var (
	Client *s3.Client
	bucket *string
)

func InitAWS() {
	awsConf := conf.Get().Aws
	if awsConf.BucketName == "" {
		return
	}
	// os.Setenv("AWS_ACCESS_KEY_ID", awsConf.AccessKey)
	// os.Setenv("AWS_SECRET_ACCESS_KEY", awsConf.SecretKey)

	cfg := aws.Config{
		Region:      awsConf.Region,
		Credentials: credentials.NewStaticCredentialsProvider(awsConf.AccessKey, awsConf.SecretKey, ""),
	}

	if awsConf.EndPoint != "" {
		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == s3.ServiceID && region == awsConf.Region {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           awsConf.EndPoint,
					SigningRegion: awsConf.Region,
				}, nil
			}
			// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})
		cfg.EndpointResolver = customResolver
	}

	Client = s3.NewFromConfig(cfg)
	bucket = aws.String(awsConf.BucketName)

	data, err := GetFile("url.txt")
	if err == nil {
		conf.SetUrl(string(*data))
		fmt.Println(conf.Get().AssetServerPath)
	}

}

func ModifiedDate(key string) (*time.Time, error) {
	if Client == nil {
		return nil, nil
	}
	data, err := Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: bucket,
		Key:    aws.String(conf.Get().Aws.BaseDir + key),
	})
	if err != nil {
		return nil, err
	}
	return data.LastModified, nil
}

func PutFile(key string, data io.Reader) error {
	if Client == nil {
		return nil
	}

	_, err := Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: bucket,
		Key:    aws.String(conf.Get().Aws.BaseDir + key),
		Body:   data,
	})
	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func GetFile(key string) (*[]byte, error) {
	if Client == nil {
		return nil, nil
	}

	data, err := Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: bucket,
		Key:    aws.String(conf.Get().Aws.BaseDir + key),
	})

	if err != nil {
		return nil, err
	}

	defer data.Body.Close()

	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

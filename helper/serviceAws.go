package helper

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/furqonzt99/snackbox/constants"
)

func UploadObjectS3(fileName string, data multipart.File) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(constants.S3_REGION),
		Credentials: credentials.NewStaticCredentials(constants.AWS_ACCESS_KEY_ID, constants.AWS_ACCESS_SECRET_KEY, ""),
	})
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(constants.S3_BUCKET),
		Key:    aws.String(fileName),
		Body:   data,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetObjectS3(fileName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(constants.S3_REGION),
		Credentials: credentials.NewStaticCredentials(constants.AWS_ACCESS_KEY_ID, constants.AWS_ACCESS_SECRET_KEY, ""),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	_, err = svc.GetObjectAcl(&s3.GetObjectAclInput{Bucket: aws.String(constants.S3_BUCKET), Key: aws.String(fileName)})
	if err != nil {
		return err
	}

	return nil
}

func DeleteObjectS3(fileName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(constants.S3_REGION),
		Credentials: credentials.NewStaticCredentials(constants.AWS_ACCESS_KEY_ID, constants.AWS_ACCESS_SECRET_KEY, ""),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(constants.S3_BUCKET), Key: aws.String(fileName)})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(constants.S3_BUCKET),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return err
	}

	return nil
}

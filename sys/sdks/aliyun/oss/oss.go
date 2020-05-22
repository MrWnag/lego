package oss

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func newOSS(opt ...Option) (IOSS, error) {
	opts := newOptions(opt...)
	if opts.Endpoint == "" || opts.AccessKeyId == "" || opts.AccessKeySecret == "" || opts.BucketName == "" {
		return nil, fmt.Errorf("参数异常 %v", opts)
	} else {
		oss := &OSS{Endpoint: opts.Endpoint, AccessKeyId: opts.AccessKeyId, AccessKeySecret: opts.AccessKeySecret, BucketName: opts.BucketName}
		err := oss.Init()
		return oss, err
	}
}

type OSS struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	client          *oss.Client
	bucket          *oss.Bucket
}

func (this *OSS) Init() (err error) {
	this.client, err = oss.New(this.Endpoint, this.AccessKeyId, this.AccessKeySecret)
	if err != nil {
		return err
	}
	if ok, err := this.client.IsBucketExist(this.BucketName); !ok || err != nil {
		this.CreateBucket(this.BucketName)
	} else {
		this.bucket, err = this.client.Bucket(this.BucketName)
	}
	return err
}

//创建存储空间。
func (this *OSS) CreateBucket(bucketName string) (err error) {
	err = this.client.CreateBucket(bucketName)
	return err
}

//上传文件
// <objectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
// <localFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
// 上传文件。
func (this *OSS) UploadFile(objectName string, localFileName string) (err error) {
	err = this.bucket.PutObjectFromFile(objectName, localFileName)
	return err
}

// 下载文件。
// <objectName>从OSS下载文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
func (this *OSS) DownloadFile(objectName string, downloadedFileName string) (err error) {
	err = this.bucket.GetObjectToFile(objectName, downloadedFileName)
	return err
}

//删除文件
func (this *OSS) DeleteFile(objectName string) (err error) {
	// 删除文件。
	err = this.bucket.DeleteObject(objectName)
	return
}

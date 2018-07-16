package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)


func (os *Ossstore) Del(path string) error {
	client, err := oss.New("Endpoint", "AccessKeyId", "AccessKeySecret")
	if err != nil {
		return err
	}

	bucket, err := client.Bucket("my-bucket")
	if err != nil {
		return err
	}

	return bucket.DeleteObject(path)
}
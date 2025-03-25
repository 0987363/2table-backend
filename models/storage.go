package models

import (
	"context"
	"errors"

	"gocloud.dev/blob"
	//	_ "gocloud.dev/blob/awsblob"
	// _ "gocloud.dev/blob/fileblob"
)

// 存储类型常量
const (
	StorageTypeLocal = "local"
	StorageTypeS3    = "s3"
)

type StorageManager struct {
	*blob.Bucket
}

// 存储配置结构体
type StorageConfig struct {
	Type  string       `json:"type" yaml:"type"`
	Local *LocalConfig `json:"local,omitempty" yaml:"local,omitempty"`
	S3    *S3Config    `json:"s3,omitempty" yaml:"s3,omitempty"`
}

type LocalConfig struct {
	Path string `json:"path" yaml:"path"`
}

type S3Config struct {
	Bucket   string `json:"bucket" yaml:"bucket"`
	Region   string `json:"region" yaml:"region"`
	Endpoint string `json:"endpoint,omitempty" yaml:"endpoint,omitempty"` // 兼容S3服务地址
}

func NewStorageManager(ctx context.Context, cfg *StorageConfig) (*StorageManager, error) {
	var bucketURL string
	var err error

	switch cfg.Type {
	case StorageTypeLocal:
		bucketURL = "file://" + cfg.Local.Path
	case StorageTypeS3:
		bucketURL = buildS3URL(cfg.S3)
	default:
		return nil, errors.New("unsupported storage type")
	}

	bucket, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		return nil, err
	}

	return &StorageManager{bucket}, nil
}

func buildS3URL(cfg *S3Config) string {
	url := "s3://" + cfg.Bucket + "?"
	if cfg.Region != "" {
		url += "region=" + cfg.Region + "&"
	}
	if cfg.Endpoint != "" {
		url += "endpoint=" + cfg.Endpoint + "&"
	}
	return url[:len(url)-1] // 去除最后一个&
}

package task

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/qq51529210/util"
)

var (
	_minio minioUploader
)

type minioUploader struct {
	// 客户端
	client *minio.Client
	// 超时
	timeout time.Duration
	// 桶
	bucket string
}

func (m *minioUploader) init(cfg *MinioCfg) error {
	m.timeout = time.Duration(cfg.Timeout) * time.Second
	m.bucket = cfg.Bucket
	//
	var err error
	// 创建客户端
	m.client, err = minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.ID, cfg.Secret, ""),
		Secure: true,
	})
	if err != nil {
		return err
	}
	// 创建桶
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	err = m.client.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}
	//
	return nil
}

func (m *minioUploader) upload(path string) (string, error) {
	// 超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	// 上传
	info, err := m.client.FPutObject(ctx, m.bucket, util.SHA1String(path), path, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return info.Key, nil
}

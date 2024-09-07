package file

import (
	"context"
	"path/filepath"
	"time"

	"github.com/lingfliu/ucs_core/ulog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioCli struct {
	EndPoint        string
	Bucket          string
	Ssl             bool
	AccessKeyId     string
	AccessKeySecret string

	minCli *minio.Client
}

func (cli *MinioCli) Open() {
	minCli, err := minio.New(cli.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cli.AccessKeyId, cli.AccessKeySecret, ""),
		Secure: cli.Ssl,
	})

	if err != nil {
		ulog.Log().E("minio", "failed to connect minio")
	}

	cli.minCli = minCli
}

func (cli *MinioCli) Close() {
	//minio client 不需要关闭
}

/**
 * @brief
 * 获取文件共享链接, 有效期72小时
 * @param fname 文件名
 */
func (cli *MinioCli) Query(fname string) string {
	url, err := cli.minCli.PresignedGetObject(context.TODO(), cli.Bucket, fname, time.Hour*72, nil)
	if err != nil {
		ulog.Log().E("minio", "failed to get object")
		return ""
	}
	return url.String()
}

func (cli *MinioCli) Upload(path string) {
	_, err := cli.minCli.FPutObject(context.Background(), cli.Bucket, filepath.Base(path), path, minio.PutObjectOptions{})
	if err != nil {
		ulog.Log().E("minio", "failed to upload object")
	}
}

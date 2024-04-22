package storage

import (
	lokiv1 "github.com/grafana/loki/operator/apis/loki/v1"
)

// Options is used to configure Loki to integrate with
// supported object storages.
type Options struct {
	Schemas     []lokiv1.ObjectStorageSchema
	SharedStore lokiv1.ObjectStorageSecretType

	Azure *AzureStorageConfig
	GCS   *GCSStorageConfig
	S3    *S3StorageConfig
	Swift *SwiftStorageConfig

	SecretName string
	SecretSHA1 string
	TLS        *TLSConfig
}

// AzureStorageConfig for Azure storage config
type AzureStorageConfig struct {
	Env       string
	Container string
}

// GCSStorageConfig for GCS storage config
type GCSStorageConfig struct {
	Bucket string
}

// S3StorageConfig for S3 storage config
type S3StorageConfig struct {
	Endpoint       string
	Region         string
	Buckets        string
	ForcePathStyle bool
}

// SwiftStorageConfig for Swift storage config
type SwiftStorageConfig struct {
	AuthURL           string
	UserDomainName    string
	UserDomainID      string
	UserID            string
	DomainID          string
	DomainName        string
	ProjectID         string
	ProjectName       string
	ProjectDomainID   string
	ProjectDomainName string
	Region            string
	Container         string
}

// TLSConfig for object storage endpoints. Currently supported only by:
// - S3
type TLSConfig struct {
	CA  string
	Key string
}

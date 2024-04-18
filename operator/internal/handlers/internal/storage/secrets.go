package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	lokiv1 "github.com/grafana/loki/operator/apis/loki/v1"
	"github.com/grafana/loki/operator/internal/manifests/storage"

	corev1 "k8s.io/api/core/v1"
)

var (
	errS3EndpointUnparseable       = errors.New("can not parse S3 endpoint as URL")
	errS3EndpointNoURL             = errors.New("endpoint for S3 must be an HTTP or HTTPS URL")
	errS3EndpointUnsupportedScheme = errors.New("scheme of S3 endpoint URL is unsupported")
	errS3EndpointAWSInvalid        = errors.New("endpoint for AWS S3 must include correct region")
)

const awsEndpointSuffix = ".amazonaws.com"

// ExtractSecret reads a k8s secret into a manifest object storage struct if valid.
func ExtractSecret(s *corev1.Secret, secretType lokiv1.ObjectStorageSecretType) (*storage.Options, error) {
	var err error
	storageOpts := storage.Options{
		SecretName:  s.Name,
		SecretSHA1:  hash,
		SharedStore: secretType,
	}

	switch secretType {
	case lokiv1.ObjectStorageSecretAzure:
		storageOpts.Azure, err = extractAzureConfigSecret(s)
	case lokiv1.ObjectStorageSecretGCS:
		storageOpts.GCS, err = extractGCSConfigSecret(s)
	case lokiv1.ObjectStorageSecretS3:
		storageOpts.S3, err = extractS3ConfigSecret(s)
	case lokiv1.ObjectStorageSecretSwift:
		storageOpts.Swift, err = extractSwiftConfigSecret(s)
	default:
		return storage.Options{}, fmt.Errorf("%w: %s", errSecretUnknownType, secretType)
	}

	if err != nil {
		return storage.Options{}, err
	}
	return storageOpts, nil
}

func hashSecretData(s *corev1.Secret) (string, error) {
	keys := make([]string, 0, len(s.Data))
	for k := range s.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha1.New()
	for _, k := range keys {
		if _, err := h.Write([]byte(k)); err != nil {
			return "", err
		}

		if _, err := h.Write(hashSeparator); err != nil {
			return "", err
		}

		if _, err := h.Write(s.Data[k]); err != nil {
			return "", err
		}

		if _, err := h.Write(hashSeparator); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func extractAzureConfigSecret(s *corev1.Secret) (*storage.AzureStorageConfig, error) {
	// Extract and validate mandatory fields
	env := s.Data[storage.KeyAzureEnvironmentName]
	if len(env) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeyAzureEnvironmentName)
	}
	container := s.Data[storage.KeyAzureStorageContainerName]
	if len(container) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeyAzureStorageContainerName)
	}
	name := s.Data[storage.KeyAzureStorageAccountName]
	if len(name) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeyAzureStorageAccountName)
	}
	key := s.Data[storage.KeyAzureStorageAccountKey]
	if len(key) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeyAzureStorageAccountKey)
	}

	return &storage.AzureStorageConfig{
		Env:       string(env),
		Container: string(container),
	}, nil
}

func extractGCSConfigSecret(s *corev1.Secret) (*storage.GCSStorageConfig, error) {
	// Extract and validate mandatory fields
	bucket := s.Data[storage.KeyGCPStorageBucketName]
	if len(bucket) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeyGCPStorageBucketName)
	}

	// Check if google authentication credentials is provided
	keyJSON := s.Data[storage.KeyGCPServiceAccountKeyFilename]
	if len(keyJSON) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeyGCPServiceAccountKeyFilename)
	}

	return &storage.GCSStorageConfig{
		Bucket: string(bucket),
	}, nil
}

func extractS3ConfigSecret(s *corev1.Secret) (*storage.S3StorageConfig, error) {
	// Extract and validate mandatory fields
	endpoint := s.Data["endpoint"]
	region := s.Data["region"]
	if err := validateS3Endpoint(string(endpoint), string(region)); err != nil {
		return nil, err
	}
	buckets := s.Data["bucketnames"]
	if len(buckets) == 0 {
		return nil, kverrors.New("missing secret field: bucketnames")
	}
	id := s.Data[storage.KeyAWSAccessKeyID]
	if len(id) == 0 {
		return nil, kverrors.New("missing secret field: access_key_id")
	}
	secret := s.Data[storage.KeyAWSAccessKeySecret]
	if len(secret) == 0 {
		return nil, kverrors.New("missing secret field: access_key_secret")
	}

	return &storage.S3StorageConfig{
		Endpoint: string(endpoint),
		Buckets:  string(buckets),
		Region:   string(region),
	}, nil
}

func validateS3Endpoint(endpoint string, region string) error {
	if len(endpoint) == 0 {
		return kverrors.New("missing secret field: endpoint")
	}

	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("%w: %w", errS3EndpointUnparseable, err)
	}

	if parsedURL.Scheme == "" {
		// Assume "just a hostname" when scheme is empty and produce a clearer error message
		return errS3EndpointNoURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: %s", errS3EndpointUnsupportedScheme, parsedURL.Scheme)
	}

	if strings.HasSuffix(endpoint, awsEndpointSuffix) {
		if len(region) == 0 {
			return kverrors.New("missing secret field: region")
		}

		validEndpointSuffix := fmt.Sprintf("%s://s3.%s%s", parsedURL.Scheme, region, awsEndpointSuffix)
		if !strings.HasSuffix(endpoint, validEndpointSuffix) {
			return fmt.Errorf("%w: %s", errS3EndpointAWSInvalid, validEndpointSuffix)
		}
	}
	return nil
}

func extractSwiftConfigSecret(s *corev1.Secret) (*storage.SwiftStorageConfig, error) {
	// Extract and validate mandatory fields
	url := s.Data[storage.KeySwiftAuthURL]
	if len(url) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeySwiftAuthURL)
	}
	username := s.Data[storage.KeySwiftUsername]
	if len(username) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeySwiftUsername)
	}
	userDomainName := s.Data[storage.KeySwiftUserDomainName]
	if len(userDomainName) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeySwiftUserDomainName)
	}
	userDomainID := s.Data[storage.KeySwiftUserDomainID]
	if len(userDomainID) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeySwiftUserDomainID)
	}
	userID := s.Data[storage.KeySwiftUserID]
	if len(userID) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeySwiftUserID)
	}
	password := s.Data[storage.KeySwiftPassword]
	if len(password) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeySwiftPassword)
	}
	domainID := s.Data[storage.KeySwiftDomainID]
	if len(domainID) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeySwiftDomainID)
	}
	domainName := s.Data[storage.KeySwiftDomainName]
	if len(domainName) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeySwiftDomainName)
	}
	containerName := s.Data[storage.KeySwiftContainerName]
	if len(containerName) == 0 {
		return nil, fmt.Errorf("%w: %s", errSecretMissingField, storage.KeySwiftContainerName)
	}

	// Extract and validate optional fields
	projectID := s.Data[storage.KeySwiftProjectID]
	projectName := s.Data[storage.KeySwiftProjectName]
	projectDomainID := s.Data[storage.KeySwiftProjectDomainId]
	projectDomainName := s.Data[storage.KeySwiftProjectDomainName]
	region := s.Data[storage.KeySwiftRegion]

	return &storage.SwiftStorageConfig{
		AuthURL:           string(url),
		UserDomainName:    string(userDomainName),
		UserDomainID:      string(userDomainID),
		UserID:            string(userID),
		DomainID:          string(domainID),
		DomainName:        string(domainName),
		ProjectID:         string(projectID),
		ProjectName:       string(projectName),
		ProjectDomainID:   string(projectDomainID),
		ProjectDomainName: string(projectDomainName),
		Region:            string(region),
		Container:         string(containerName),
	}, nil
}

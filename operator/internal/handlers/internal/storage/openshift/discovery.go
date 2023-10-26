package openshift

import "os"

const (
	identityTokenPath = "/var/run/secrets/openshift/serviceaccount/token"
)

type AWSSTSEnv struct {
	RoleARN              string
	WebIdentityTokenFile string
}

type AzureWIFEnv struct {
	ClientID       string
	TenantID       string
	SubscriptionID string
}

type STSEnv struct {
	AWS   *AWSSTSEnv
	Azure *AzureWIFEnv
}

func DiscoverSTSEnv() *STSEnv {
	var (
		// AWS
		roleARN = os.Getenv("ROLEARN")
		// Azure
		clientID       = os.Getenv("CLIENTID")
		tenantID       = os.Getenv("TENANTID")
		subscriptionID = os.Getenv("SUBSCRIPTIONID")
	)

	switch {
	case roleARN != "":
		return &STSEnv{
			AWS: &AWSSTSEnv{
				RoleARN:              roleARN,
				WebIdentityTokenFile: identityTokenPath,
			},
		}
	case clientID != "" && tenantID != "" && subscriptionID != "":
		return &STSEnv{
			Azure: &AzureWIFEnv{
				ClientID:       clientID,
				TenantID:       tenantID,
				SubscriptionID: subscriptionID,
			},
		}
	}

	return &STSEnv{}
}

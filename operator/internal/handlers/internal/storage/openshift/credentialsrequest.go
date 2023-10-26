package openshift

import (
	"context"
	"fmt"

	"github.com/ViaQ/logerr/v2/kverrors"
	cloudcredentialv1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/grafana/loki/operator/internal/external/k8s"
)

const (
	ccoNamespace = "openshift-cloud-credential-operator"
)

func CreateCredentialsRequest(ctx context.Context, k k8s.Client, stack client.ObjectKey, sts *STSEnv) (client.ObjectKey, error) {
	providerSpec, name, tokenPath, err := encodeProvideSpec(stack.Name, sts)

	credReqKey := client.ObjectKey{Name: name, Namespace: stack.Namespace}

	if err != nil {
		return credReqKey, kverrors.Wrap(err, "failed encoding provider spec")
	}

	credReq := &cloudcredentialv1.CredentialsRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ccoNamespace,
		},
		Spec: cloudcredentialv1.CredentialsRequestSpec{
			SecretRef: corev1.ObjectReference{
				Name:      name,
				Namespace: stack.Namespace,
			},
			ProviderSpec: providerSpec,
			ServiceAccountNames: []string{
				stack.Name,
			},
			CloudTokenPath: tokenPath,
		},
	}

	if err := k.Create(ctx, credReq); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return credReqKey, kverrors.Wrap(err, "failed to create credentialsrequest", "key", client.ObjectKeyFromObject(credReq))
		}
	}

	return credReqKey, nil
}

func encodeProvideSpec(stackName string, sts *STSEnv) (*runtime.RawExtension, string, string, error) {
	var (
		spec                   runtime.Object
		credReqName, tokenPath string
	)

	switch {
	case sts.AWS != nil:
		spec = &cloudcredentialv1.AWSProviderSpec{
			StatementEntries: []cloudcredentialv1.StatementEntry{
				{
					Action: []string{
						"s3:ListBucket",
						"s3:PutObject",
						"s3:GetObject",
						"s3:DeleteObject",
					},
					Effect:   "Allow",
					Resource: "arn:aws:s3:*:*:*",
				},
			},
			STSIAMRoleARN: sts.AWS.RoleARN,
		}
		credReqName = fmt.Sprintf("%s-aws-creds", stackName)
		tokenPath = sts.AWS.WebIdentityTokenFile
	case sts.Azure != nil:
		// TODO(@periklis) Missing implementation from sts.Azure to AzureProviderSpec
		//                 Requires OCP 4.15 or final backport to 4.14.z!!!
		//                 Waiting for: https://github.com/openshift/cloud-credential-operator/pull/587
		//                         and: https://github.com/openshift/console/pull/13082
		spec = &cloudcredentialv1.AzureProviderSpec{}
		credReqName = fmt.Sprintf("%s-azure-creds", stackName)
		tokenPath = identityTokenPath
	}

	encodedSpec, err := cloudcredentialv1.Codec.EncodeProviderSpec(spec.DeepCopyObject())
	return encodedSpec, credReqName, tokenPath, err
}

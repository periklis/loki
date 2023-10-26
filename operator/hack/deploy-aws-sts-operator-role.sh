#!/bin/bash

#
# AWS definitions
#
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
OIDC_PROVIDER=$(oc get authentication cluster -ojson | jq -r .spec.serviceAccountIssuer | sed -e "s/^https:\/\///")

#
# Product definitions
#
PRODUCT_NAME=${1:-"openshift-logging"}
NAMESPACE=${2:-"openshift-operators-redhat"}
SERVICE_ACCOUNT_NAME=${3:-"loki-operator-controller-manager"}
OPERAND_NAMESPACE=${4:-"openshift-logging"}
OPERAND_SERVICE_ACCOUNT_NAME=${5:-"lokistack-dev"}

#
# STS Role and Policy definitions
#
ROLE_NAME="$PRODUCT_NAME-$SERVICE_ACCOUNT_NAME"
POLICY_ARN="arn:aws:iam::aws:policy/AmazonS3FullAccess"

read -r -d '' TRUST_RELATIONSHIP <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Effect": "Allow",
     "Principal": {
       "Federated": "arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/${OIDC_PROVIDER}"
     },
     "Action": "sts:AssumeRoleWithWebIdentity",
     "Condition": {
       "StringEquals": {
         "${OIDC_PROVIDER}:sub": [
           "system:serviceaccount:${NAMESPACE}:${SERVICE_ACCOUNT_NAME}",
           "system:serviceaccount:${OPERAND_NAMESPACE}:${OPERAND_SERVICE_ACCOUNT_NAME}"
         ]
       }
     }
   }
 ]
}
EOF

echo "${TRUST_RELATIONSHIP}" > /tmp/trust.json

export AWS_PAGER=""
ROLE_ARN=$(aws iam create-role --role-name "$ROLE_NAME" --assume-role-policy-document file:///tmp/trust.json --description "AWS STS Loki Operator Role for $PRODUCT_NAME" --query Role.Arn --output text)

echo -n "Attaching $POLICY_ARN ... "
aws iam attach-role-policy \
    --role-name "$ROLE_NAME" \
    --policy-arn "${POLICY_ARN}"
echo "ok."

echo -n "Use $ROLE_ARN when installing your operator on the STS cluster."

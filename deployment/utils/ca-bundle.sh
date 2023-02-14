#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

export CA_BUNDLE=$(kubectl get configmap -n kube-system extension-apiserver-authentication -o=jsonpath='{.data.client-ca-file}' | base64 |tr -d '\n')
cat > $(dirname $(dirname "$(realpath -s "$0")"))/values/ca-bundle.yaml <<EOF
webhook:
  caBundle: "${CA_BUNDLE}"
EOF

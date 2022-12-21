#!/bin/bash

export GOPATH="$(go env GOPATH)"

# pkg/client/clientset
# pkg/versioned/typed

input=${GROUPSVERSION:-cloud/v1alpha1}

client-gen --clientset-name clientset \
--input-base 'github.com/apecloud/kubeblocks/apis' \
--input $input \
--output-base "${GOPATH}/src/" \
--output-package 'github.com/apecloud/kubeblocks/pkg' \
--go-header-file ./hack/boilerplate.go.txt \
-v=9
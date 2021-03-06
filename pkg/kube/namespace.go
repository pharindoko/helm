/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kube // import "k8s.io/helm/pkg/kube"

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

func createNamespace(client internalclientset.Interface, namespace string) error {
	ns := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := client.Core().Namespaces().Create(ns)
	return err
}

func getNamespace(client internalclientset.Interface, namespace string) (*core.Namespace, error) {
	return client.Core().Namespaces().Get(namespace, metav1.GetOptions{})
}

func ensureNamespace(client internalclientset.Interface, namespace string) error {
	_, err := getNamespace(client, namespace)
	if err != nil && errors.IsNotFound(err) {
		err = createNamespace(client, namespace)

		// If multiple commands which run `ensureNamespace` are run in
		// parallel, then protect against the race condition in which
		// the namespace did not exist when `getNamespace` was executed,
		// but did exist when `createNamespace` was executed. If that
		// happens, we can just proceed as normal.
		if errors.IsAlreadyExists(err) {
			return nil
		}
	}
	return err
}

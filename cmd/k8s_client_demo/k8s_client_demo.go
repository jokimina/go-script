/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"reflect"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func getResouceNames(inter interface{}) (nameList []string) {
	l := reflect.Indirect(reflect.ValueOf(inter)).FieldByName("Items")
	switch l.Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(l.Interface())
		for i := 0; i < s.Len(); i++ {
			item := reflect.ValueOf(s.Index(i).Interface())
			name := item.FieldByName("Name").Interface().(string)
			nameList = append(nameList, name)
		}
	default:
		fmt.Println(l.Kind())
		fmt.Println(reflect.ValueOf(l))
	}
	return
}

func int32Ptr(i int32) *int32 { return &i }

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d namespaces in the cluster, %v\n", len(namespaces.Items), getResouceNames(namespaces))

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	pods, err = clientset.CoreV1().Pods(metav1.NamespaceDefault).List(metav1.ListOptions{})
	fmt.Println(getResouceNames(pods))

	// deployment
	deployClient := clientset.AppsV1().Deployments(metav1.NamespaceDefault)
	deployments, err := deployClient.List(metav1.ListOptions{LabelSelector: "env=testing"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("There are %d deployment on %s namespace, in the cluster\n", len(deployments.Items), metav1.NamespaceDefault)
	fmt.Println(getResouceNames(deployments))

	for _, deployment := range deployments.Items {
		appContainer := deployment.Spec.Template.Spec.Containers[0]
		//fmt.Println(deploy.String())
		fmt.Println(deployment.GetName(), appContainer.Image)
	}

	deployment, err := deployClient.Get("nginx", metav1.GetOptions{})
	if err != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err))
	}
	appContainer := deployment.Spec.Template.Spec.Containers[0]
	deployment.Spec.Replicas = int32Ptr(2)
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Set(int64(2000000))
	deployment.Spec.Template.Spec.Containers[0].Image = "nginx:alpine"
	fmt.Println(deployment.GetName(), appContainer.Image)
	_, err = deployClient.UpdateStatus(deployment)
	b, err := json.MarshalIndent(deployment.Spec.Template.Spec.Containers, "", "  ")
	fmt.Println(string(b))
	if err != nil {
		panic(err)
	}
}

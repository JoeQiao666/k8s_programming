package main

import (
	"cloudstate-client-go/pkg/statefulservice"
	"cloudstate-client-go/pkg/statefulstore"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	cloudstate "github.com/cloudstateio/cloudstate/cloudstate-operator/pkg/apis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func createStatefulStore(kubeconfig string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	statefulstoreClient, err := statefulstore.NewClient(config)
	if err != nil {
		panic(err)
	}

	statefulstore := &cloudstate.StatefulStore{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   "statefulstore-wei",
			Labels: map[string]string{"mylabel": "stafulstorecrd"},
		},
		Spec: cloudstate.StatefulStoreSpec{
			InMemory: true,
		},
		Status: cloudstate.StatefulStoreStatus{
			Summary: "created",
		},
	}
	// Create the statefulstore object we create above in the k8s cluster
	ctx := context.Background()
	resp, err := statefulstoreClient.Statefulstore("default").Create(statefulstore, ctx)
	if err != nil {
		fmt.Printf("error while creating object: %v\n", err)
	} else {
		fmt.Printf("object created: %v\n", resp)
	}

	obj, err := statefulstoreClient.Statefulstore("default").Get(statefulstore.ObjectMeta.Name, ctx)
	if err != nil {
		fmt.Printf("error while getting the object %v\n", err)
	}
	fmt.Printf("statefulstore Objects Found: \n%+v\n", obj)
}

func createStatefulService(kubeconfig string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	statefulserviceClient, err := statefulservice.NewClient(config)
	if err != nil {
		panic(err)
	}
	hostEnv := corev1.EnvVar{
		Name:  "HOST",
		Value: "localhost",
	}
	portEnv := corev1.EnvVar{
		Name:  "PORT",
		Value: "8080",
	}
	container := corev1.Container{
		Name:  "credit-user-function",
		Image: "gcr.io/sap-nj-serverless-poc/credit:latest",
		Env: []corev1.EnvVar{
			hostEnv,
			portEnv,
		},
	}
	statefulservice := &cloudstate.StatefulService{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   "statefulservice-wei",
			Labels: map[string]string{"mylabel": "stafulservicecrd"},
		},
		Spec: cloudstate.StatefulServiceSpec{
			Containers: []corev1.Container{container},
			StoreConfig: &cloudstate.StatefulServiceStoreConfig{
				// Database: "inmemory",
				StatefulStore: corev1.LocalObjectReference{
					Name: "statefulstore-wei",
				},
			},
		},
		Status: cloudstate.StatefulServiceStatus{
			Summary:  "created",
			Replicas: 2,
		},
	}
	// Create the statefulservice object we create above in the k8s cluster
	ctx := context.Background()
	resp, err := statefulserviceClient.Statefulservice("default").Create(statefulservice, ctx)
	if err != nil {
		fmt.Printf("error while creating object: %v\n", err)
	} else {
		fmt.Printf("object created: %v\n", resp)
	}

	obj, err := statefulserviceClient.Statefulservice("default").Get(statefulservice.ObjectMeta.Name, ctx)
	if err != nil {
		fmt.Printf("error while getting the object %v\n", err)
	}
	fmt.Printf("statefulservice Objects Found: \n%+v\n", obj)
}

func getEndpoints(kubeconfig string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	eps, err := clientset.CoreV1().Endpoints("default").List(ctx, meta_v1.ListOptions{FieldSelector: "metadata.name=statefulservice-wei"})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%v\n", eps)
}

func main() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	createStatefulStore(kubeconfig)
	fmt.Println()
	createStatefulService(kubeconfig)
	fmt.Println()
	getEndpoints(kubeconfig)
}

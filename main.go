package main

import (
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/alighm/sample-controller/controllers"
	"github.com/alighm/sample-controller/pkg/client/clientset/versioned"
	"github.com/alighm/sample-controller/pkg/client/informers/externalversions"
)

func main() {
	log.Info("Sample-Controller")

	var (
		config *rest.Config
		err    error
	)

	kubeconfig := ""
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "kubeconfig file")
	flag.Parse()
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		log.Errorf("error creating client: %v", err)
		os.Exit(1)
	}

	// get kube client for connectivity
	client, crClient := getKubernetesClient(config)

	// retrieve our custom resource informer which was generated from the code generator and pass it the custom
	// resource client, specifying we should be looking through all namespaces for listing and watching
	informer := externalversions.NewSharedInformerFactory(crClient, 10*time.Minute)

	// creating the controller and passing it
	controller := controllers.NewController(client, informer.Foo().V1().HelloTypes())

	informer.Start(nil)
	controller.Run(nil)
}

// retrieve the Kubernetes cluster client from outside of the cluster
func getKubernetesClient(config *rest.Config) (kubernetes.Interface, versioned.Interface) {
	// generate the client based off of the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	myresourceClient, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Info("Successfully constructed k8s client")
	return client, myresourceClient
}

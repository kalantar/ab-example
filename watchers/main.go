package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	watchtool "k8s.io/client-go/tools/watch"
)

var namespaces []string = []string{"default"}
var resources []schema.GroupVersionResource = []schema.GroupVersionResource{{
	Group:    "apps",
	Version:  "v1",
	Resource: "deployments",
}, {
	Group:    "",
	Version:  "v1",
	Resource: "services",
}}

// channel used to terminate go processes
var stopCh chan struct{} = make(chan struct{})

func main() {
	config, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	client, _ := dynamic.NewForConfig(config)

	// The iter8-abn service maintains a single list of applications. As applications/versions are identified,
	// this list is updated. If it is possible for simultaneous events be processed at the same time, mutual
	// exclusion on the shared application structure will be required. Alternatively, we can order updates to
	// avoid this requirement.

	// We demonstrate three solutions:

	// (a) dynamic.ReourceInterface.Watch() - this is the basic watch api for dynamic resources
	// To avoid races, the events from each watcher are sent to a common channel for processing
	go demonstrateWatchInterface(client)

	// (b) watchtools.RetryWatch - this wraps the above to handle unexpected failures of the watchers
	// To avoid races, the events from each watcher are sent to a common channel for processing
	go demonstrateRetryWatch(client)

	// (c) informer - this is a lower? level API that does not have issue with failure and already orders events
	go demonstrateInformer(client)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)

	<-sigCh
	close(stopCh)
}

func demonstrateWatchInterface(client dynamic.Interface) {
	// channel used to order events from all objects
	var allEvents chan watch.Event = make(chan watch.Event)

	// for each namespace, resource type configure watch
	for _, ns := range namespaces {
		for _, gvr := range resources {
			w, _ := client.Resource(gvr).Namespace(ns).Watch(context.Background(), metav1.ListOptions{})
			go func() {
				for e := range w.ResultChan() {
					allEvents <- e
				}
			}()
		}
	}

	// forever process ordered list of events
	for {
		event := <-allEvents
		switch event.Type {
		case watch.Added:
			Add("Watch.Interface", event.Object.(*unstructured.Unstructured))
		case watch.Modified:
			Update("Watch.Interface", event.Object.(*unstructured.Unstructured))
		case watch.Deleted:
			Delete("Watch.Interface", event.Object.(*unstructured.Unstructured))
		}
	}
}

func demonstrateRetryWatch(client dynamic.Interface) {
	// channel used to order events from all objects
	var allEvents chan watch.Event = make(chan watch.Event)

	// for each namespace, resource configure RetryWatch
	for _, ns := range namespaces {
		for _, gvr := range resources {
			watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
				return client.Resource(gvr).Namespace(ns).Watch(context.Background(), metav1.ListOptions{})
			}
			w, _ := watchtool.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})
			go func() {
				for e := range w.ResultChan() {
					allEvents <- e
				}
			}()
		}
	}

	// forever process ordered list of events
	for {
		event := <-allEvents
		switch event.Type {
		case watch.Added:
			Add("RetryWatch", event.Object.(*unstructured.Unstructured))
		case watch.Modified:
			Update("RetryWatch", event.Object.(*unstructured.Unstructured))
		case watch.Deleted:
			Delete("RetryWatch", event.Object.(*unstructured.Unstructured))
		}
	}
}

func demonstrateInformer(client dynamic.Interface) {
	// channel used to terminate informers
	stopCh := make(chan struct{})

	// for each namespace, resource type configure Informer
	for _, ns := range namespaces {
		for _, gvr := range resources {
			factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, 0, ns, nil)
			informer := factory.ForResource(gvr)
			informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc:    func(obj interface{}) { Add("Informer", obj.(*unstructured.Unstructured)) },
				UpdateFunc: func(oldObj, obj interface{}) { Update("Informer", obj.(*unstructured.Unstructured)) },
				DeleteFunc: func(obj interface{}) { Delete("Informer", obj.(*unstructured.Unstructured)) },
			})
			factory.Start(stopCh)
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)

	<-sigCh
	close(stopCh)
}

func Add(label string, obj *unstructured.Unstructured) {
	fmt.Printf("%s ADD called for %s/%s\n", label, obj.GetNamespace(), obj.GetName())
}

func Update(label string, obj *unstructured.Unstructured) {
	fmt.Printf("%s UPDATE called for %s/%s\n", label, obj.GetNamespace(), obj.GetName())
}

func Delete(label string, obj *unstructured.Unstructured) {
	fmt.Printf("%s DELETE called for %s/%s\n", label, obj.GetNamespace(), obj.GetName())
}

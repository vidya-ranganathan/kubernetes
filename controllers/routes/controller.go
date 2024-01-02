package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset      kubernetes.Interface
	depLister      appslisters.DeploymentLister
	depCacheSynced cache.InformerSynced
	queue          workqueue.RateLimitingInterface
}

func newController(clientset kubernetes.Interface, depInformer appsinformers.DeploymentInformer) *controller {
	con := &controller{
		clientset:      clientset,
		depLister:      depInformer.Lister(),
		depCacheSynced: depInformer.Informer().HasSynced,
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "routes"),
	}

	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    con.handleAdd,
			DeleteFunc: con.handleDelete,
		},
	)

	return con
}

func (con *controller) run(ch <-chan struct{}) {
	fmt.Println("starting routes controller")
	if !cache.WaitForCacheSync(ch, con.depCacheSynced) {
		fmt.Println("waiting for cache to be synced")
	}
	go wait.Until(con.worker, 1*time.Second, ch)

	<-ch
}

func (con *controller) worker() {
	/*
		for con.processItem() {

		}
	*/
}

func (con *controller) handleAdd(obj interface{}) {
	fmt.Println("add invoked..")
	con.queue.Add(obj)
}

func (con *controller) handleDelete(obj interface{}) {
	fmt.Println("delete invoked..")
	con.queue.Add(obj)
}

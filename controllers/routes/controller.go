package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	// perform the controller task here ....
	for con.controllerTask() {

	}

}

func (con *controller) controllerTask() bool {
	// get the object and if not Q is empty , return false
	item, shutdown := con.queue.Get()
	if shutdown {
		return false
	}

	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("failed getting key from cache %s\n", err.Error())
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("failed splitting key into name and namespace -  %s\n", err.Error())
	}

	err = con.syncDeployment(ns, name)
	if err != nil {
		fmt.Printf("failed to sync deployment - %s", err.Error())
		return false
	}
	return true
}

func (con *controller) syncDeployment(ns, name string) error {
	ctx := context.Background()

	// creating service for the deployment
	svc := corev1.Service{}
	_, err := con.clientset.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error while creating service for the deployment %s\n", err.Error())
	}

	return nil
}

func (con *controller) handleAdd(obj interface{}) {
	fmt.Println("add invoked..")
	con.queue.Add(obj)
}

func (con *controller) handleDelete(obj interface{}) {
	fmt.Println("delete invoked..")
	con.queue.Add(obj)
}

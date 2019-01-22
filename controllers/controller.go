package controllers

import (
	"time"

	log "github.com/sirupsen/logrus"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/alighm/sample-controller/pkg/client/informers/externalversions/foo/v1"
)

type Controller struct {
	helloTypeListerSynced cache.InformerSynced
	queue                 workqueue.RateLimitingInterface
}

func NewController(client kubernetes.Interface, helloTypeInformer v1.HelloTypeInformer) *Controller {
	c := &Controller{
		helloTypeListerSynced: helloTypeInformer.Informer().HasSynced,
		queue:                 workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "hellotype"),
	}

	helloTypeInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				log.Info("HelloType Added")
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				log.Info("HelloType Updated")
			},
			DeleteFunc: func(obj interface{}) {
				log.Info("HelloType Deleted")
			},
		},
	)

	return c
}

// Run will start the controller.
// StopCh channel is used to send interrupt signal to stop it.
func (c *Controller) Run(stopCh <-chan struct{}) {
	// don't let panics crash the process
	defer utilruntime.HandleCrash()

	// make sure the work queue is shutdown which will trigger workers to end
	defer c.queue.ShutDown()

	// wait for the caches to synchronize before starting the worker
	if !cache.WaitForCacheSync(stopCh, c.helloTypeListerSynced) {
		log.Error("timed out waiting for caches to sync")
		return
	}

	// runWorker will loop until "something bad" happens.  The .Until will
	// then rekick the worker after one second
	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) runWorker() {
	// processNextWorkItem will automatically wait until there's work available
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	return true
}
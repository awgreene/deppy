package entitysource

import (
	"github.com/go-logr/logr"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

type DeppyRegistryV1Cache struct {
	sync.Mutex
	watch  watch.Interface
	client client.Client
	logger *logr.Logger
	cache  map[string][]EntityID
}

func NewDeppyRegistryV1Cache(watch watch.Interface, client client.Client, logger *logr.Logger) *DeppyRegistryV1Cache {
	return &DeppyRegistryV1Cache{
		watch:  watch,
		client: client,
		logger: logger,
	}
}

func (r *DeppyRegistryV1Cache) StartCache() {
	for {
		select {
		case entry := <-r.watch.ResultChan():
			r.logger.V(1).Info("New catalogSource event", "catalogsource", entry.Object, "event", entry.Type)
			switch entry.Type {
			case watch.Deleted:
				r.Mutex.Lock()
				defer r.Mutex.Unlock()
				delete(r.cache)
			case watch.Added:
			case watch.Modified:
			}
		}
	}
}

func namespacedName(o *client.Object) string {
	catsrc, ok := o.(*v1alpha1.CatalogSource)
	if !ok {
		return ""
	}
	return types.NamespacedName{Namespace: catsrc.Namespace, Name: catsrc.Name}.String()
}

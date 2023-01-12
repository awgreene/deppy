package entitysource

import (
	"sync"

	"github.com/go-logr/logr"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
			cs, ok := entry.Object.(*v1alpha1.CatalogSource)
			if !ok {
				r.logger.V(1).Info(typecastErrorMessage)
				return
			}

			switch entry.Type {
			case watch.Deleted:
				r.Mutex.Lock()
				defer r.Mutex.Unlock()
				delete(r.cache, namespacedName(cs))
			case watch.Added, watch.Modified:
				r.Mutex.Lock()
				defer r.Mutex.Unlock()
				needsUpdate, err := r.cacheNeedsUpdate(cs)
				if err != nil {
					r.logger.V(1).Info(typecastErrorMessage)
					return
				}

				if !needsUpdate {
					r.logger.V(1).Info("Cache is up to date", "CatalogSource", namespacedName(cs))
					return
				}

				// Update cache

			}
		}
	}
}

func namespacedName(cs *v1alpha1.CatalogSource) string {
	return types.NamespacedName{Namespace: cs.GetNamespace(), Name: cs.GetName()}.String()
}

const (
	typecastErrorMessage = "unable to convert event object into a catalogSource"
)

func (r *DeppyRegistryV1Cache) cacheNeedsUpdate(cs *v1alpha1.CatalogSource) (bool, error) {
	return true, nil
}

func (r *DeppyRegistryV1Cache) retrieveBundles(cs *v1alpha1.CatalogSource) ([]Entity, error) {
	svc := cs.Spec.Address

	// use svc to retrieve entities
	return nil, nil
}

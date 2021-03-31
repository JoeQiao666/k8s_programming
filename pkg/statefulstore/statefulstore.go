package statefulstore

import (
	"context"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	cloudstate "github.com/cloudstateio/cloudstate/cloudstate-operator/pkg/apis/v1alpha1"
)

type StatefulstoreClient struct {
	restClient rest.Interface
}

func (c *StatefulstoreClient) Statefulstore(namespace string) StatefulstoreInterface {
	return &statefulstoreClient{
		client: c.restClient,
		ns:     namespace,
	}
}

type StatefulstoreInterface interface {
	Create(obj *cloudstate.StatefulStore, ctx context.Context) (*cloudstate.StatefulStore, error)
	Update(obj *cloudstate.StatefulStore, ctx context.Context) (*cloudstate.StatefulStore, error)
	Delete(name string, options *meta_v1.DeleteOptions, ctx context.Context) error
	Get(name string, ctx context.Context) (*cloudstate.StatefulStore, error)
}

type statefulstoreClient struct {
	client rest.Interface
	ns     string
}

func (c *statefulstoreClient) Create(obj *cloudstate.StatefulStore, ctx context.Context) (*cloudstate.StatefulStore, error) {
	result := &cloudstate.StatefulStore{}
	err := c.client.Post().
		Namespace(c.ns).Resource("statefulstores").
		Body(obj).Do(ctx).Into(result)
	return result, err
}

func (c *statefulstoreClient) Update(obj *cloudstate.StatefulStore, ctx context.Context) (*cloudstate.StatefulStore, error) {
	result := &cloudstate.StatefulStore{}
	err := c.client.Put().
		Namespace(c.ns).Resource("statefulstores").
		Body(obj).Do(ctx).Into(result)
	return result, err
}

func (c *statefulstoreClient) Delete(name string, options *meta_v1.DeleteOptions, ctx context.Context) error {
	return c.client.Delete().
		Namespace(c.ns).Resource("statefulstores").
		Name(name).Body(options).Do(ctx).
		Error()
}

func (c *statefulstoreClient) Get(name string, ctx context.Context) (*cloudstate.StatefulStore, error) {
	result := &cloudstate.StatefulStore{}
	err := c.client.Get().
		Namespace(c.ns).Resource("statefulstores").
		Name(name).Do(ctx).Into(result)
	return result, err
}

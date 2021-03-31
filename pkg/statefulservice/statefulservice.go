package statefulservice

import (
	"context"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	cloudstate "github.com/cloudstateio/cloudstate/cloudstate-operator/pkg/apis/v1alpha1"
)

type StatefulserviceClient struct {
	restClient rest.Interface
}

func (c *StatefulserviceClient) Statefulservice(namespace string) StatefulserviceInterface {
	return &statefulserviceClient{
		client: c.restClient,
		ns:     namespace,
	}
}

type StatefulserviceInterface interface {
	Create(obj *cloudstate.StatefulService, ctx context.Context) (*cloudstate.StatefulService, error)
	Update(obj *cloudstate.StatefulService, ctx context.Context) (*cloudstate.StatefulService, error)
	Delete(name string, options *meta_v1.DeleteOptions, ctx context.Context) error
	Get(name string, ctx context.Context) (*cloudstate.StatefulService, error)
}

type statefulserviceClient struct {
	client rest.Interface
	ns     string
}

func (c *statefulserviceClient) Create(obj *cloudstate.StatefulService, ctx context.Context) (*cloudstate.StatefulService, error) {
	result := &cloudstate.StatefulService{}
	err := c.client.Post().
		Namespace(c.ns).Resource("statefulservices").
		Body(obj).Do(ctx).Into(result)
	return result, err
}

func (c *statefulserviceClient) Update(obj *cloudstate.StatefulService, ctx context.Context) (*cloudstate.StatefulService, error) {
	result := &cloudstate.StatefulService{}
	err := c.client.Put().
		Namespace(c.ns).Resource("statefulservices").
		Body(obj).Do(ctx).Into(result)
	return result, err
}

func (c *statefulserviceClient) Delete(name string, options *meta_v1.DeleteOptions, ctx context.Context) error {
	return c.client.Delete().
		Namespace(c.ns).Resource("statefulservices").
		Name(name).Body(options).Do(ctx).
		Error()
}

func (c *statefulserviceClient) Get(name string, ctx context.Context) (*cloudstate.StatefulService, error) {
	result := &cloudstate.StatefulService{}
	err := c.client.Get().
		Namespace(c.ns).Resource("statefulservices").
		Name(name).Do(ctx).Into(result)
	return result, err
}

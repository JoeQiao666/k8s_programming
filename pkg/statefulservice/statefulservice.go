package statefulservice

import (
	"context"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
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
	Watch(ctx context.Context, opts meta_v1.ListOptions) (watch.Interface, error)
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

func (c *statefulserviceClient) Watch(ctx context.Context, opts meta_v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	scheme := runtime.NewScheme()
	schemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := schemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return c.client.Get().
		Namespace(c.ns).
		Resource("statefulservices").
		VersionedParams(&opts, runtime.NewParameterCodec(scheme)).
		Timeout(timeout).
		Watch(ctx)
}

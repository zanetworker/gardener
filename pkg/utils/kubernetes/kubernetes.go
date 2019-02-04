package kubernetes

import (
	"context"
	"fmt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SetMetaDataLabel sets the key value pair in the labels section of the given ObjectMeta.
// If the given ObjectMeta did not yet have labels, they are initialized.
func SetMetaDataLabel(meta *metav1.ObjectMeta, key, value string) {
	if meta.Labels == nil {
		meta.Labels = make(map[string]string)
	}

	meta.Labels[key] = value
}

// HasMetaDataAnnotation checks if the passed meta object has the given key, value set in the annotations section.
func HasMetaDataAnnotation(meta *metav1.ObjectMeta, key, value string) bool {
	val, ok := meta.Annotations[key]
	return ok && val == value
}

// CreateOrUpdate creates or updates the object. Optionally, it executes a transformation function before the
// request is made.
func CreateOrUpdate(ctx context.Context, c client.Client, obj runtime.Object, transform func() error) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	if err := c.Get(ctx, key, obj); err != nil {
		if apierrors.IsNotFound(err) {
			if transform != nil && transform() != nil {
				return err
			}
			return c.Create(ctx, obj)
		}
		return err
	}

	if transform != nil && transform() != nil {
		return err
	}
	return c.Update(ctx, obj)
}

func Key(namespaceOrName string, nameOpt ...string) client.ObjectKey {
	if len(nameOpt) == 0 {
		return client.ObjectKey{Name: namespaceOrName}
	}
	if len(nameOpt) > 1 {
		panic(fmt.Sprintf("more than name/namespace for key specified: %s/%v", namespaceOrName, nameOpt))
	}
	return client.ObjectKey{Namespace: namespaceOrName, Name: nameOpt[0]}
}

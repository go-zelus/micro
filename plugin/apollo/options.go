package apollo

import (
	"context"
	"github.com/micro/go-micro/v2/config/source"
)

type addressKey struct{}
type namespaceKey struct{}
type clusterKey struct{}
type appIdKey struct{}
type backupConfigPathKey struct{}

func WithAddress(address string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, addressKey{}, address)
	}
}

func WithNamespace(namespace string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, namespaceKey{}, namespace)
	}
}

func WithCluster(cluster string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, clusterKey{}, cluster)
	}
}

func WithAppID(appId string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, appIdKey{}, appId)
	}
}

func WithBackupConfigPath(p string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, backupConfigPathKey{}, p)
	}
}

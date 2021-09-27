package apollo

import (
	"fmt"
	"strings"
	"time"

	"github.com/apolloconfig/agollo/v4"
	apolloconfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/logger"
)

type apolloSource struct {
	client    *agollo.Client
	opts      source.Options
	namespace string
}

func (a *apolloSource) Read() (*source.ChangeSet, error) {
	cache := a.client.GetConfigCache(a.namespace)

	kv := map[string]interface{}{}

	cache.Range(func(key, value interface{}) bool {
		kv[key.(string)] = string(value.([]byte))
		return true
	})

	kv = convert(kv)
	b, err := a.opts.Encoder.Encode(kv)
	if err != nil {
		return nil, fmt.Errorf("error reading source: %v", err)
	}
	scs := &source.ChangeSet{
		Data:      b,
		Format:    a.opts.Encoder.String(),
		Source:    a.String(),
		Timestamp: time.Now(),
	}
	scs.Checksum = scs.Sum()
	return scs, nil
}

func (a *apolloSource) Write(set *source.ChangeSet) error {
	return nil
}

func (a *apolloSource) Watch() (source.Watcher, error) {
	watcher := NewWatcher(a.String(), a.opts.Encoder)
	return watcher, nil
}

func (a *apolloSource) String() string {
	return "apollo"
}

func NewApolloSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	namespace, ok := options.Context.Value(namespaceKey{}).(string)
	if !ok {
		logger.Error("apollo namespace error")
	}
	address, ok := options.Context.Value(addressKey{}).(string)
	if !ok {
		logger.Error("apollo address error")
	}

	appId, ok := options.Context.Value(appIdKey{}).(string)
	if !ok {
		logger.Error("apollo appid error")
	}

	cluster, ok := options.Context.Value(clusterKey{}).(string)
	if !ok {
		cluster = "dev"
	}

	backupConfigPath, ok := options.Context.Value(backupConfigPathKey{}).(string)
	if !ok {
		backupConfigPath = "./config"
	}

	logger.Infof("[apollo config] address: %s, namespace: %s", address, namespace)

	readyConfig := &apolloconfig.AppConfig{
		AppID:            appId,
		Cluster:          cluster,
		IP:               address,
		NamespaceName:    namespace,
		IsBackupConfig:   true,
		BackupConfigPath: backupConfigPath,
	}

	client, err := agollo.StartWithConfig(func() (*apolloconfig.AppConfig, error) {
		return readyConfig, nil
	})

	if err != nil {
		logger.Fatal(err)
	}

	return &apolloSource{
		client:    client,
		opts:      options,
		namespace: namespace,
	}
}

func convert(kv map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for k, v := range kv {
		if k == "" {
			continue
		}

		target := data
		sp := strings.Split(k, ".")
		for _, dir := range sp[:len(sp)-1] {
			if _, ok := target[dir]; !ok {
				target[dir] = make(map[string]interface{})
			}
			target = target[dir].(map[string]interface{})
		}
		leaf := sp[len(sp)-1]
		target[leaf] = v
	}
	return data
}

package apollo

import (
	"github.com/go-zelus/micro/config"
	"github.com/go-zelus/micro/plugin/apollo"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/encoder/yaml"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/logger"
)

func NewSource(c *cli.Context) source.Source {
	address := c.String("apollo_address")
	if len(address) == 0 {
		address = config.Env("APOLLO_ADDRESS")
	}
	if len(address) == 0 {
		logger.Fatal("need config apollo_address")
	}

	namespace := c.String("apollo_namespace")
	if len(namespace) == 0 {
		namespace = config.Env("APOLLO_NAMESPACE", "application")
	}

	appId := c.String("apollo_app_id")
	if len(appId) == 0 {
		appId = config.Env("APOLLO_APPID")
	}

	cluster := c.String("apollo_cluster")
	if len(cluster) == 0 {
		cluster = config.Env("APOLLO_CLUSTER", "dev")
	}

	backupConfigPath := config.Env("BACKUP_CONFIG_PATH", "./")

	coder := yaml.NewEncoder()
	return apollo.NewApolloSource(
		apollo.WithAppID(appId),
		apollo.WithAddress(address),
		apollo.WithNamespace(namespace),
		apollo.WithCluster(cluster),
		apollo.WithBackupConfigPath(backupConfigPath),
		source.WithEncoder(coder),
	)
}

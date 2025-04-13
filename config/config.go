package config

type TracesConfig struct {
	Endpoint string `env:"CONFIG__TRACES_CONFIG__ENDPOINT" default:"http://tempo.observability.svc.cluster.local:14268"`
	Path     string `env:"CONFIG__TRACES_CONFIG__PATH" default:"/api/traces"`
}

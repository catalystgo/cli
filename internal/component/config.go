package component

var _ Component = configComponent{}

type (
	configComponent struct{}
)

// Content implements Component.
func (c configComponent) Content() ([]byte, error) {
	return []byte(`
# General configuration
app:
  name: "catalyst_go_app"

# Server configuration
server:
  http:
    enabled: true
    port: 8082
  grpc:
    port: 8081
  tls:
    enabled: false
    cert_file: "/path/to/cert.pem"
    key_file: "/path/to/key.pem"

# Logging configuration
logging:
  level: "info"
  format: "json"

# Tracing configuration
tracing:
  enabled: true
  provider: "jaeger"  # or "zipkin", "opentelemetry", etc.
  address: "localhost:6831"

# Vault configuration
vault:
  enable: true
  address: "localhost:8200"

realtime_config: # Registered in ETCD
   -  name: log_level # ["ERROR", "WARN", "INFO", "DEBUG"]
      usage: Log level enum
      value: "INFO"
      type: string

# Other custom configurations
custom:
  key1: "value1"
  key2: "value2"
	condition: service_started

`), nil
}

// Name implements Component.
func (c configComponent) Name() string {
	return "config.yml"
}

// Path implements Component.
func (c configComponent) Path() string {
	return "."
}

func NewConfigComponent() Component {
	return configComponent{}
}

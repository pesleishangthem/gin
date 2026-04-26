package config

// ServerConfig is what the Gin module needs to know to run.
type ServerConfig interface {
	GetPort() string
	GetMode() string             // "release" or "debug"
	GetKeycloakRealmURL() string // realm URL
}

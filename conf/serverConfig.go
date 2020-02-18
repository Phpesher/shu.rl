package conf

type ServerConfig struct {
	ServerPort string
	ServerHost string
}

func NewConfig(port, host string) *ServerConfig {
	return &ServerConfig{":8080", "localhost"}
}
package conf

// Config is application config
type Config struct {
	Server   Server   `json:"server" yaml:"server"`
	Database Database `json:"database" yaml:"database"`
	Logging  Logging  `json:"logging" yaml:"logging"`
}

type Server struct {
	Port string `json:"port" yaml:"port"`
	Cert string `json:"cert" yaml:"cert"`
	Key  string `json:"cert" yaml:"key"`
	TLS  bool   `json:"tls" yaml:"tls"`
}

type Database struct {
	Type         string `json:"type" yaml:"type"`
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	User         string `json:"user" yaml:"user"`
	Password     string `json:"password" yaml:"password"`
	DatabaseName string `json:"databaseName" yaml:"databaseName"`
	SslMode      string `json:"sslMode" yaml:"sslMode"`
	SslFactory   string `json:"sslFactory" yaml:"sslFactory"`
}

type Logging struct {
	Level string `json:"level" yaml:"level"`
}

// SaneDefaults provides base config for testing
func SaneDefaults() Config {
	var config = Config{
		Server: Server{
			Port: "8090",
			Cert: "certs/cert.crt",
			Key:  "certs/cert.key",
			TLS:  false,
		},
		Database: Database{
			Host:         "127.0.0.1",
			Port:         5432,
			User:         "user",
			Password:     "password",
			DatabaseName: "test",
			SslMode:      "disable",
			SslFactory:   "org.postgresql.ssl.NonValidatingFactory",
		},
		Logging: Logging{
			Level: "debug",
		},
	}
	return config
}

package config

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func NewConfig() *Config {
	config := &Config{}

	config.Username = "postgres"
	config.Password = "postgres"
	config.Host = "localhost"
	config.Port = "5432"
	config.Database = "db_tasks"

	return config
}

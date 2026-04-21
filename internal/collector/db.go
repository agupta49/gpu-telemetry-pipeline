package collector

import "fmt"

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewConfig(host, port, user string) *Config {
	return &Config{Host: host, Port: port, User: user}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}

func (c *Config) Validate() error {
	if c.Host == "" || c.Port == "" || c.User == "" {
		return fmt.Errorf("missing required config fields")
	}
	return nil
}

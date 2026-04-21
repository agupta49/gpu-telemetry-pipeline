package collector

import "fmt"

type Config struct {
	PostgresHost string
	PostgresPort int
	PostgresUser string
	PostgresDB   string
	MQAddress    string
}

func NewConfig(host, user, db string) *Config {
	return &Config{
		PostgresHost: host,
		PostgresPort: 5432,
		PostgresUser: user,
		PostgresDB: db,
		MQAddress: "mq:50051",
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", 
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresDB)
}

func (c *Config) Validate() error {
	if c.PostgresHost == "" {
		return fmt.Errorf("postgres host required")
	}
	if c.PostgresDB == "" {
		return fmt.Errorf("postgres db required")
	}
	return nil
}

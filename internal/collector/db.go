package collector
type Config struct{ DSN string }
func NewConfig() *Config { return &Config{} }
func (c *Config) DSNStr() string { return c.DSN }
func (c *Config) Validate() error { return nil }

package config

type Config struct {
	Port int
	Env  string
	Db   struct {
		DBname       string
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
}

type Option func(с *Config)

func NewConfig(opts ...Option) *Config {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}

	if config.Port == 0 {
		config.Port = 8080
	}

	if config.Env == "" {
		config.Env = "development"
	}

	if config.Db.DBname == "" {
		config.Db.DBname = "test"
	}

	if config.Db.Dsn == "" {
		config.Db.Dsn = "test"
	}

	if config.Db.MaxOpenConns == 0 {
		config.Db.MaxOpenConns = 50
	}

	if config.Db.MaxIdleConns == 0 {
		config.Db.MaxIdleConns = 50
	}

	if config.Db.MaxIdleTime == "" {
		config.Db.MaxIdleTime = "15m"
	}
	return config
}

func WithPort(port int) Option {
	return func(c *Config) { c.Port = port }
}

func WithEnv(env string) Option {
	return func(c *Config) { c.Env = env }
}

func WithDBname(db string) Option {
	return func(c *Config) { c.Db.DBname = db }
}

func WithDSN(dsn string) Option {
	return func(c *Config) { c.Db.Dsn = dsn }
}

func WithMaxOpenConns(maxOpenConns int) Option {
	return func(c *Config) { c.Db.MaxOpenConns = maxOpenConns }
}

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(c *Config) { c.Db.MaxIdleConns = maxIdleConns }
}

func WithMaxIdleTime(maxIdleTime string) Option {
	return func(c *Config) { c.Db.MaxIdleTime = maxIdleTime }
}

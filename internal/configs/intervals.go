package configs

type ClientIntervalsCfg struct {
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	HashKey        string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

type ServerIntervalsCfg struct {
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	HashKey         string `env:"KEY"`
}

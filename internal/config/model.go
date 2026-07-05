package config

type Config struct {
	AppName     string           `json:"app_name"`
	Environment string           `json:"environment"`
	HTTP        HTTPConfig       `json:"http"`
	Databases   []DatabaseConfig `json:"databases"`
	Log         LogConfig        `json:"log"`
}

type HTTPConfig struct {
	ExternalPort string `json:"external_port"`
	InternalPort string `json:"internal_port"`
}

type DatabaseConfig struct {
	Nickname           string    `json:"nickname"`
	Host               string    `json:"host"`
	Port               int       `json:"port"`
	User               string    `json:"user"`
	Password           string    `json:"password"`
	Name               string    `json:"name"`
	SSLMode            string    `json:"ssl_mode"`
	SSLConfig          SSLConfig `json:"ssl_config"`
	MaxOpenConns       int       `json:"max_open_conns"`
	MaxIdleConns       int       `json:"max_idle_conns"`
	TransactionTimeout int       `json:"transaction_timeout"`
	IsDefault          bool      `json:"is_default"`
	ReadOnly           bool      `json:"read_only"`
}

type SSLConfig struct {
	ClientCert string `json:"client_cert"`
	ClientKey  string `json:"client_key"`
	RootCA     string `json:"root_ca"`
}

type LogConfig struct {
	Level            string `json:"level"`
	Encoding         string `json:"encoding"`
	Development      bool   `json:"development"`
	OutputPaths      string `json:"output_paths"`
	ErrorOutputPaths string `json:"error_output_paths"`
}

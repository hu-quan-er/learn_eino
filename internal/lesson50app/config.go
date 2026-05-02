package lesson50app

type Config struct {
	AppName string
	Tenant  string
}

type Result struct {
	Final string
	Trace []string
}

func DefaultConfig() Config {
	return Config{
		AppName: "lesson50_cli_app",
		Tenant:  "tutorial-team",
	}
}

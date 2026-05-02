package lesson46app

type Config struct {
	AgentName string
	Prefix    string
}

func DefaultConfig() Config {
	return Config{
		AgentName: "lesson46_bootstrap_agent",
		Prefix:    "bootstrap",
	}
}

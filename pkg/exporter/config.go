package exporter

// Config holds the configuration for the ethereum sync status tool.
type Config struct {
	// Execution is the execution node to use.
	Execution ExecutionNode `yaml:"execution"`
	// ConsensusNodes is the consensus node to use.
	Consensus ConsensusNode `yaml:"consensus"`
	// PollingFrequencySeconds determines how frequently to poll the targets if running in Serve mode.
	PollingFrequencySeconds int `yaml:"pollingFrequencySeconds"`
	// DiskUsage determines if the disk usage metrics should be exported.
	DiskUsage DiskUsage `yaml:"diskUsage"`
}

// ConsensusNode represents a single ethereum consensus client.
type ConsensusNode struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
}

// ExecutionNode represents a single ethereum execution client.
type ExecutionNode struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
}

// DiskUsage configures the exporter to expose disk usage stats for these directories.
type DiskUsage struct {
	Enabled     bool     `yaml:"enabled"`
	Directories []string `yaml:"directories"`
}

// DefaultConfig represents a sane-default configuration.
func DefaultConfig() *Config {
	return &Config{
		Execution: ExecutionNode{
			Enabled: true,
			Name:    "execution",
			URL:     "http://localhost:8545",
		},
		Consensus: ConsensusNode{
			Enabled: true,
			Name:    "consensus",
			URL:     "http://localhost:5052",
		},
		PollingFrequencySeconds: 5,
		DiskUsage: DiskUsage{
			Enabled:     false,
			Directories: []string{},
		},
	}
}
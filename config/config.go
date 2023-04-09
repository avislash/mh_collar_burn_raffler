package config

import "time"

type Config struct {
	MaxWinners         uint64           `yaml:"max_winners"`
	MetadataEndpoint   string           `yaml:"metadata_endpoint"`
	EtherscanRateLimit int              `yaml:"etherscan_rate_limit"`
	Snapshot           Snapshot         `yaml:"snapshot"`
	IneligibleTraits   IneligibleTraits `yaml:"ineligible_traits"`
}

type Snapshot struct {
	Start time.Time
	Stop  time.Time
}
type IneligibleTraits struct {
	Forms  []string
	Faces  []string
	Mouths []string
	Torsos []string
}

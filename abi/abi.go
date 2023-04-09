package abi

import _ "embed"

var (
	//go:embed mh_collars_proxy.json
	MutantHoundCollarsProxyABI string

	//go:embed mh_collars.json
	MutantHoundCollarsABI string
)

package main

type Host struct {
	FQDN string `yaml:"fqdn"`
	Role string `yaml:"role"`
}

func (h *Host) getHost() string {
	if h.FQDN == "" {
		return "localhost"
	}
	return h.FQDN
}

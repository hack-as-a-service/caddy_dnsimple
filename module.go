package caddy_dnsimple

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/hack-as-a-service/caddy_dnsimple/provider"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct{ *provider.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.dnsimple",
		New: func() caddy.Module { return &Provider{new(provider.Provider)} },
	}
}

// TODO: This is just an example. Useful to allow env variable placeholders; update accordingly.
// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.Provider.APIToken = caddy.NewReplacer().ReplaceAll(p.Provider.APIToken, "")
	p.Provider.AccountID = caddy.NewReplacer().ReplaceAll(p.Provider.AccountID, "")
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*Provider)(nil)
)

package provider

// Provider is struct contains all contract of provider
// use when you have outbound rest http transaction
type Provider struct {
}

// NewProvider to initiate provider
// if you have another dependencies, please add on Deps struct, not on the params
func NewProvider(dep *Deps) *Provider {
	return &Provider{}
}

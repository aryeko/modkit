package module

// Token identifies a provider for resolution.
type Token string

// Resolver provides access to resolved provider instances.
type Resolver interface {
	Get(Token) (any, error)
}

// ResolverFunc adapts a function to a Resolver.
type ResolverFunc func(Token) (any, error)

// Get implements Resolver.
func (f ResolverFunc) Get(token Token) (any, error) {
	return f(token)
}

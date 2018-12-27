package adapters

// CacheAdapter - Adapter to talk to cache
type ESLAdapter interface {
	Originate(key string) (string, error)
	GetVar(key string) (string)
}

package adapters

// CacheAdapter - Adapter to talk to cache
type ESLAdapter interface {
	SendBgApiCmd(key string) (string, error)
	SendApiCmd(key string) (string)
}

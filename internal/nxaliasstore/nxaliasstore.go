package nxaliasstore

// NxAliasStore is a storage for nxalias aliass
type NxAliasStore struct {
	NxAlias map[string]string
}

// New creates a new storage for nxalias aliass
func New() *NxAliasStore {
	return &NxAliasStore{
		NxAlias: make(map[string]string),
	}
}

// New creates a new alias in the map
func (ws *NxAliasStore) New(domain string, alias string) {
	ws.NxAlias[domain] = alias
}

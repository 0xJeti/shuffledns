package wildcardstore

// WildcardStore is a storage for wildcard domains
type WildcardStore struct {
	Domain map[string]struct{}
}

// New creates a new storage for wildcard domains
func New() *WildcardStore {
	return &WildcardStore{
		Domain: make(map[string]struct{}),
	}
}

// New creates a new domain in the map
func (ws *WildcardStore) New(domain string) {
	ws.Domain[domain] = struct{}{}
}

// Exists indicates if an domain exists in the map
func (ws *WildcardStore) Exists(domain string) bool {
	_, ok := ws.Domain[domain]
	return ok
}

// Get gets the meta-information for an IP address from the map.
func (ws *WildcardStore) Get(domain string) {
}

// Delete deletes the records for an IP from store.
func (ws *WildcardStore) Delete(domain string) {
	delete(ws.Domain, domain)
}

// Close removes all the references to arrays and releases memory to the gc
func (ws *WildcardStore) Close() {
}

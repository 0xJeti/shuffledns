package rcodestore

// RCodeStore is a storage for response codes
type RCodeStore struct {
	RCode map[string]*RCodeMeta
}

// RCodeMeta contains meta-information about single
// error response code found during enumeration.
type RCodeMeta struct {
	RCode    string
	Resolver string
}

// New creates a new storage for response codes
func New() *RCodeStore {
	return &RCodeStore{
		RCode: make(map[string]*RCodeMeta),
	}
}

// New creates a new response code in the map
func (rs *RCodeStore) New(domain, rcode, resolver string) {
	rs.RCode[domain] = &RCodeMeta{RCode: rcode, Resolver: resolver}
}

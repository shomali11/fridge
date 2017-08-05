package fridge

// RetrievalOption an option for retrieval
type RetrievalOption func(*RetrievalDetails)

// WithRestock sets retrieval restocking option
func WithRestock(restock func() (string, error)) RetrievalOption {
	return func(retrievalInfo *RetrievalDetails) {
		retrievalInfo.Restock = restock
	}
}

// RetrievalDetails contains retrieval information
type RetrievalDetails struct {
	Restock func() (string, error)
}

func newRetrievalDetails(options ...RetrievalOption) *RetrievalDetails {
	retrievalDetails := &RetrievalDetails{}
	for _, option := range options {
		option(retrievalDetails)
	}
	return retrievalDetails
}

package api

type EtlApi struct {
	globalIndex int
	dataRows []StatementRow
}

func NewEtlApi() *EtlApi {
	return &EtlApi{
		globalIndex: 0,
		dataRows:    make([]StatementRow, 0),
	}
}

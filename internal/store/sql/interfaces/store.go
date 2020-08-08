package interfaces

type (
	SqlStoreDriver string

	SqlTxStore interface {
		Resource() Resource
		ResourceResourceRel() ResourceResourceRel
	}
)

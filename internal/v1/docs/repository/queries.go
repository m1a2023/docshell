package repository

const (
	get_all_documents  = "select * from documents;"
	get_document_by_id = "select * from documents where id = $1;"
)

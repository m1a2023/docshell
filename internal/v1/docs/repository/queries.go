package repository

const (
	get_all_documents  = "select * from documents;"
	get_document_by_id = "select * from documents where id = $1;"
	insert_document    = `
		insert into documents (
			author_id, uploader_id, title, size, path, hash
		)
			values (
				$1, $2, $3, $4, $5, $6
				) returning *;
	`
)

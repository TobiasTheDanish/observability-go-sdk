package store

type Metadata struct {
	Id           string
	ExportPath   string
	ExportMethod string
}

type Data interface {
	Metadata() Metadata
}

type Store interface {
	/*Store a new chunk of data*/
	Store(data Data)
	/*Get all data from store WITHOUT invalidating*/
	GetAll() []Data
	/*Invalidate data from store, based on the provided metadata id*/
	Invalidate(id string)
}

package api

// Serializer converts models to/from JSON representations
type Serializer[T any] interface {
	ToRepresentation(obj *T) (interface{}, error)
	FromCreate(body []byte) (*T, error)
	FromUpdate(obj *T, body []byte) error
}


package codec

type Encoder interface {
	Encode(v interface{}) error
}
type Decoder interface {
	Decode(v interface{}) error
	More() bool
}

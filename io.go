package puresound

type Reader func(buf []byte) (int, error)

func (r Reader) Read(buf []byte) (int, error) {
	return r(buf)
}


type Writer func(buf []byte) (int, error)

func (w Writer) Write(buf []byte) (int, error) {
	return w(buf)
}

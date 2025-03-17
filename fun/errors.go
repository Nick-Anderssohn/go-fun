package fun

type errorConst string

const errEndOfStream errorConst = "end of stream reached"

func (e errorConst) Error() string {
	return string(e)
}

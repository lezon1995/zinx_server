package util

func NewBuffer() []byte {
	return make([]byte, GlobalConfig.BufferSize)
}
func NewLenBuffer(len uint32) []byte {
	return make([]byte, len)
}

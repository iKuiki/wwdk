package storer

// Storer 运行时存储器
type Storer interface {
	// WriterString 刷写配置
	WriterString(data string) error
	// ReadString 读取配置
	ReadString() (data string, err error)
}

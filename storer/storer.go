package storer

// Storer 运行时存储器
type Storer interface {
	// Writer 刷写配置
	Writer(data []byte) error
	// Read 读取配置
	Read() (data []byte, err error)
	// Truncate 清空配置
	Truncate() (err error)
}

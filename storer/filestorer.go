package storer

import (
	"io/ioutil"
	"os"
)

// FileStorer 文件配置存储器
type FileStorer struct {
	file *os.File
}

// MustNewFileStorer 创建文件配置存储器
func MustNewFileStorer(filePath string) Storer {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	return &FileStorer{
		file: file,
	}
}

// WriterString 刷写配置
func (storer *FileStorer) WriterString(data string) error {
	err := storer.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = storer.file.WriteString(data)
	if err != nil {
		return err
	}
	// seek file point to start
	_, err = storer.file.Seek(0, 0)
	return err
}

// ReadString 读取配置
func (storer *FileStorer) ReadString() (data string, err error) {
	d, err := ioutil.ReadAll(storer.file)
	if err != nil {
		return "", err
	}
	// seek file point to start
	_, err = storer.file.Seek(0, 0)
	return string(d), err
}

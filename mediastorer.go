package wwdk

import (
	"os"

	"github.com/pkg/errors"
)

// MediaType 媒体类型
type MediaType int32

const (
	// MediaTypeUserHeadImg 用户头像媒体类型
	MediaTypeUserHeadImg MediaType = 1
	// MediaTypeContactHeadImg 联系人头像媒体类型
	MediaTypeContactHeadImg MediaType = 2
	// MediaTypeMessageImage 信息图片媒体类型
	MediaTypeMessageImage MediaType = 3
	// MediaTypeMessageVoice 信息音频媒体类型
	MediaTypeMessageVoice MediaType = 4
	// MediaTypeMessageVideo 信息视频媒体类型
	MediaTypeMessageVideo MediaType = 5
)

// MediaFile 媒体文件
type MediaFile struct {
	// MediaType 媒体类型
	MediaType MediaType
	// FileName 文件名
	FileName string
	// BinaryContent 文件的二进制内容
	BinaryContent []byte
}

// MediaStorer 媒体文件储存器
type MediaStorer interface {
	// Storer 储存媒体文件，传入媒体文件，返回媒体文件URL与err异常
	Storer(file MediaFile) (url string, err error)
}

// localMediaStorer 内置的媒体储存器，将媒体文件储存到当前文件夹下
type localMediaStorer struct {
	saveDir string
}

// NewLocalMediaStorer 新建本地媒体存储器
func NewLocalMediaStorer(saveDir string) MediaStorer {
	return &localMediaStorer{
		saveDir: saveDir,
	}
}

// Storer 储存媒体文件
func (s *localMediaStorer) Storer(file MediaFile) (url string, err error) {
	filename := s.saveDir + file.FileName
	f, err := os.Create(filename)
	if err != nil {
		return "", errors.New("create " + filename + " error: " + err.Error())
	}
	defer f.Close()
	_, err = f.Write(file.BinaryContent)
	if err != nil {
		return "", errors.New("Write to file error: " + err.Error())
	}
	return filename, nil
}

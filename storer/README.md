# Storer 配置存储器

---

storer是用于存储运行配置的存储器，可以提供存、取一串string配置，只能存取一个字符串数据


调用```Truncate() (err error)```清空数据
调用```WriterString(data string) error```方法可以写入数据
调用```ReadString() (data string, err error)```则读出数据

数据的parse与format需要自行处理

fileStorre提供了一个storer的基于文件的实现
需要用MustNewFileStorer方法来初始化

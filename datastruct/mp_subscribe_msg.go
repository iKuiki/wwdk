package datastruct

// MPSubscribeMsg 订阅的公众号消息
type MPSubscribeMsg struct {
	MPArticleCount int64       `json:"MPArticleCount"`
	MPArticleList  []MPArticle `json:"MPArticleList"`
	NickName       string      `json:"NickName"`
	Time           int64       `json:"Time"`
	UserName       string      `json:"UserName"`
}

// MPArticle 公众号的文章
type MPArticle struct {
	Cover  string `json:"Cover"`
	Digest string `json:"Digest"`
	Title  string `json:"Title"`
	URL    string `json:"Url"`
}

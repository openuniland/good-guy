package types

type ArticleCrawl struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Aid   int    `json:"aid"`
}

type UpdatedWithNewArticleResponse struct {
	Data          []*ArticleCrawl `json:"data"`
	IsNew         bool            `json:"is_new"`
	SubscribedIDs []string        `bson:"subscribed_ids" json:"subscribed_ids"`
}

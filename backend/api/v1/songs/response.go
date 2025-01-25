package songs

type urlResponse struct {
	Url string `json:"url"`
}

func newUrlResponse(url string) *urlResponse {
	r := &urlResponse{}
	r.Url = url
	return r
}

package playlist

type urlResponse struct {
	Url string `json:"url"`
}

type playlistResponse struct {
	PlaylistID   string `json:"id"`
	Name string `json:"name"`
}

type songResponse struct {
	SongID string `json:"id"`
}

func newUrlResponse(url string) *urlResponse {
	r := &urlResponse{}
	r.Url = url
	return r
}

func newPlaylistResponse(id string, name string) *playlistResponse {
	r := &playlistResponse{}
	r.PlaylistID = id
	r.Name = name
	return r
}

func newSongResponse(id string) *songResponse {
	r := &songResponse{}
	r.SongID = id
	return r
}
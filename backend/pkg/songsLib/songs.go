package songsLib

import (
	"fmt"
	"math"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"backend/models"
	"backend/store"
	"time"

	soundcloudapi "github.com/zackradisic/soundcloud-api"
)

type SongDetails struct {
	Title      string `json:"title"`
	DurationMS int64  `json:"duration_ms"`
	PlayUrl    string `json:"play_url"`
}

// better name, we need struct to have an place to put the client for soundcloud
type SongGetter struct {
	sc *soundcloudapi.API
	ss *store.QueueStore
	cs *store.CacheStore
}

func NewSongGetter(ss *store.QueueStore, cs *store.CacheStore) (*SongGetter, error) {
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})

	if err != nil {
		return nil, err
	}

	return &SongGetter{sc: sc, ss: ss, cs: cs}, nil
}

func (s *SongGetter) GetSongDetails(song_url string) (SongDetails, error) {
	u, err := url.Parse(song_url)
	if err != nil {
		return SongDetails{}, err
	}

	valid, err := s.VerifyURL(u)
	if !valid || err != nil {
		return SongDetails{}, fmt.Errorf("invalid url")
	}

	source := s.GetSongSource(u)
	song_id := s.GetSongId(u, source)

	song_in_cache, err := s.cs.IsSongInCache(song_id, source)
	if err != nil {
		return SongDetails{}, fmt.Errorf("failed to get song: %w", err)
	}

	if song_in_cache {
		song, err := s.cs.GetSong(song_id, source)
		if err != nil {
			return SongDetails{}, fmt.Errorf("failed to get song: %w", err)
		}

		fmt.Println("returning song from cache")
		if song.ExpiresAt.After(time.Now()) {
			return SongDetails{
				Title:      song.Title,
				DurationMS: int64(song.DurationMS),
				PlayUrl:    song.PlayURL,
			}, nil
		}
	}

	switch source {
		case "youtube":
			return s.GetYoutubeDetails(song_url)
		case "soundcloud":
			return s.GetSoundcloudDetails(song_url)
	}

	return SongDetails{}, fmt.Errorf("unsupported url")
}

func (s *SongGetter) GetSongSource(u *url.URL) string {
	source := ""
	if strings.HasSuffix(u.Host, "youtube.com") || strings.HasSuffix(u.Host, "youtu.be") {
		source = "youtube"
	} else if strings.HasSuffix(u.Host, "soundcloud.com") {
		source = "soundcloud"
	}
	return source
}

func (s *SongGetter) VerifyURL(song_url *url.URL) (bool, error) {
	host := song_url.Host
	if strings.HasSuffix(host, "youtube.com") ||
		strings.HasSuffix(host, "youtu.be") ||
		strings.HasSuffix(host, "soundcloud.com") {
		return true, nil
	}

	return false, nil
}

func (s *SongGetter) GetYoutubeDetails(url string) (SongDetails, error) {
	cmd := exec.Command("/app/bin/yt-dlp", "-f", "bestaudio", "--get-title", "--get-duration", "-g", url)
	output, err := cmd.Output()
	if err != nil {
		return SongDetails{}, fmt.Errorf("failed to get youtube details: %w", err)
	}

	output_str := strings.TrimSpace(string(output))
	fmt.Println(output_str)
	splited := strings.Split(output_str, "\n")

	splited_time := strings.Split(splited[2], ":")
	duration_s := int64(0)
	//TODO: test, code generated
	for i := 0; i < len(splited_time); i++ {
		val, err := strconv.ParseInt(splited_time[i], 10, 64)
		if err != nil {
			return SongDetails{}, fmt.Errorf("failed to parse duration: %w", err)
		}
		duration_s += int64(math.Pow(60, float64(len(splited_time)-i-1))) * val
	}
	duration_ms := duration_s * 1000

	details := SongDetails{
		Title:      splited[0],
		DurationMS: duration_ms,
		PlayUrl:    splited[1],
	}

	defer s.insertYoutubeSongInCache(url, details)

	return details, nil
}

func (s *SongGetter) GetSoundcloudDetails(url string) (SongDetails, error) {
	tracks, err := s.sc.GetTrackInfo(soundcloudapi.GetTrackInfoOptions{
		URL: url,
	})
	if err != nil {
		return SongDetails{}, fmt.Errorf("failed to get soundcloud details: %w", err)
	}

	if len(tracks) == 0 {
		return SongDetails{}, fmt.Errorf("no tracks found")
	}

	track := tracks[0]
	play_url, err := s.sc.GetDownloadURL(url, "progressive")
	if err != nil {
		return SongDetails{}, fmt.Errorf("failed to get download url: %w", err)
	}

	details := SongDetails{
		Title:      track.Title,
		DurationMS: track.FullDurationMS,
		PlayUrl:    play_url,
	}

	defer s.insertSoundcloudSongInCache(url, details)

	return details, nil
}

func (s *SongGetter) GetSongId(u *url.URL, source string) string {
	if source == "youtube" {
		return u.Query().Get("v")
	} else if source == "soundcloud" {
		return u.Path
	}

	return ""
}

func (s *SongGetter) insertYoutubeSongInCache(song_url string, song_details SongDetails) {
	u, err := url.Parse(song_url)
	if err != nil {
		fmt.Printf("failed to parse youtube url: %v\n", err)
		return
	}

	video_id := u.Query().Get("v")
	if video_id == "" {
		fmt.Println("failed to get video id")
		return
	}
	
	play_url, err := url.Parse(song_details.PlayUrl)
	if err != nil {
		fmt.Println("failed to parse youtube url: %w", err)
		return
	}

	expires := play_url.Query().Get("expire")
	if expires == "" {
		fmt.Println("failed to get expires")
		return
	}

	expires_at, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		fmt.Println("failed to parse expires: %w", err)
		return
	}

	s.insertSongInCache(video_id, song_url, "youtube", expires_at, song_details)
}

func (s *SongGetter) insertSoundcloudSongInCache(song_url string, song_details SongDetails) {
	u, err := url.Parse(song_url)
	if err != nil {
		fmt.Printf("failed to parse youtube url: %v\n", err)
		return
	}

	track_id := u.Path
	//TODO: test how long lings last for, currently 24 hours for testing 
	s.insertSongInCache(track_id, song_url, "soundcloud", time.Now().Add(time.Hour*24).Unix(), song_details)
}


func (s *SongGetter) insertSongInCache(song_id string, song_url string, source string, expires_at int64, song_details SongDetails) {
	err := s.cs.InsertSong(&models.SongMapping{
		SongID:     song_id,
		Source:     source,
		ExpiresAt:  time.Unix(expires_at, 0),
		SongURL:    song_url,
		Title:      song_details.Title,
		DurationMS: int(song_details.DurationMS),
		PlayURL:    song_details.PlayUrl,
	})
	if err != nil {
		fmt.Printf("failed to insert song in cache: %v\n", err)
	}
}


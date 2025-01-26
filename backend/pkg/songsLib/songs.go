package songsLib

import (
	"fmt"
	"math"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	soundcloudapi "github.com/zackradisic/soundcloud-api"
)

type SongDetails struct {
	Title      string `json:"title"`
	DurationMS int64  `json:"duration_ms"`
	PlayUrl    string `json:"play_url"`
}

// better name, we need stuckt to have an place to put the client for soundcloud
type SongGetter struct {
	sc *soundcloudapi.API
}

func NewSongGetter() (*SongGetter, error) {
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})

	if err != nil {
		return nil, err
	}

	return &SongGetter{sc: sc}, nil
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

	if strings.HasSuffix(u.Host, "youtube.com") || strings.HasSuffix(u.Host, "youtu.be") {
		return s.GetYoutubeDetails(song_url)
	} else if strings.HasSuffix(u.Host, "soundcloud.com") {
		return s.GetSoundcloudDetails(song_url)
	}

	return SongDetails{}, fmt.Errorf("unsupported url")
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
	duration_ms := int64(0)
	//TODO: test, code generated
	for i := 0; i < len(splited_time); i++ {
		val, err := strconv.ParseInt(splited_time[i], 10, 64)
		if err != nil {
			return SongDetails{}, fmt.Errorf("failed to parse duration: %w", err)
		}
		duration_ms += int64(math.Pow(60, float64(len(splited_time)-i-1))) * val
	}

	return SongDetails{
		Title:      splited[0],
		DurationMS: duration_ms,
		PlayUrl:    splited[1],
	}, nil
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
	play_url := ""
	for _, format := range track.Media.Transcodings {
		if format.Preset == "mp3_0_1" && format.Format.Protocol == "hls" { //TODO: do we want hsl ? I think so
			play_url = format.URL
			break
		}
	}

	if play_url == "" {
		return SongDetails{}, fmt.Errorf("no mp3_0_1 hsl format found")
	}

	return SongDetails{
		Title:      track.Title,
		DurationMS: track.FullDurationMS,
		PlayUrl:    play_url,
	}, nil
}

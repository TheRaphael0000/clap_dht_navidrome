package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/navidrome/navidrome/plugins/pdk/go/host"
	"github.com/navidrome/navidrome/plugins/pdk/go/metadata"
	"github.com/navidrome/navidrome/plugins/pdk/go/pdk"
	"github.com/navidrome/navidrome/plugins/pdk/go/types"
)

const configAPIUrl = "apiUrl"

type TrackResponse struct {
	Path       string  `json:"path"`
	ExternalId string  `json:"external_id"`
	Score      float64 `json:"score"`
}

type clapdhtPlugin struct{}

const LOG_PREFIX = "[ClapDHT]"

func logInfo(text string) {
	pdk.Log(pdk.LogInfo, fmt.Sprintf("%s %s", LOG_PREFIX, text))
}

func logError(text error) {
	pdk.Log(pdk.LogError, fmt.Sprintf("%s %s", LOG_PREFIX, text.Error()))
}

func init() {
	metadata.Register(&clapdhtPlugin{})
	logInfo("Plugin Loaded")
}

func (p *clapdhtPlugin) GetSimilarSongsByTrack(input metadata.SimilarSongsByTrackRequest) (*metadata.SimilarSongsResponse, error) {
	logInfo(fmt.Sprintf("GetSimilarSongsByTrack (track ID: %s, Name: %s, Artist: %s)", input.ID, input.Name, input.Artist))

	tracks, err := p.queryAPI(input.ID, int(input.Count))
	if err != nil {
		logError(err)
		return nil, err
	}

	// Convert to Navidrome SongRef format preserving order
	songs := make([]types.SongRef, 0, len(tracks))
	for _, track := range tracks {
		songs = append(songs, types.SongRef{
			ID: track.ExternalId,
		})
	}

	logInfo(fmt.Sprintf("%d songs found", len(songs)))

	return &metadata.SimilarSongsResponse{Songs: songs}, nil
}

func (p *clapdhtPlugin) queryAPI(itemID string, count int) ([]TrackResponse, error) {
	apiBaseURL, ok := pdk.GetConfig(configAPIUrl)

	if !ok {
		return nil, fmt.Errorf("No API URL set")
	}

	params := url.Values{}
	params.Set("limit", strconv.Itoa(count))

	apiURL := fmt.Sprintf("%s/query/%s?%s", apiBaseURL, itemID, params.Encode())
	logInfo(fmt.Sprintf("query API: %s", apiURL))

	resp, err := host.HTTPSend(host.HTTPRequest{
		Method: "GET",
		URL:    apiURL,
	})

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response status %d", resp.StatusCode)
	} else {
		logInfo(fmt.Sprintf("response status: %d", resp.StatusCode))
	}

	var tracks []TrackResponse
	if err := json.Unmarshal(resp.Body, &tracks); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tracks, nil
}

func main() {}

package main

import (
	"fmt"
	"net/http"
	"strings"
)

type BBCUpdate struct {
	Realtime struct {
		Artist            string  `json:"artist"`
		BrandPid          string  `json:"brand_pid"`
		End               float64 `json:"end"`
		EpisodePid        string  `json:"episode_pid"`
		MusicbrainzArtist struct {
			ID       string      `json:"id"`
			Name     string      `json:"name"`
			SortName string      `json:"sort_name"`
			Type     interface{} `json:"type"`
		} `json:"musicbrainz_artist"`
		ProgrammeOffset float64 `json:"programme_offset"`
		RecordID        string  `json:"record_id"`
		SecondsAgo      float64 `json:"seconds_ago"`
		SegmentEventPid string  `json:"segment_event_pid"`
		Start           float64 `json:"start"`
		Title           string  `json:"title"`
		Type            string  `json:"type"`
		VersionPid      string  `json:"version_pid"`
	} `json:"realtime"`
	RequestMaxSeconds float64 `json:"requestMaxSeconds"`
	RequestMinSeconds float64 `json:"requestMinSeconds"`
}

func main() {

}

func trimBBCoutput(input string) string {
	t := strings.Replace(input, "realtimeCallback(", "", 1)
	return t[:len(t)-1]
}

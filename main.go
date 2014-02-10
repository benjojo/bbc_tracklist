package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// http://polling.bbc.co.uk/radio/realtime/bbc_1xtra.jsonp
	fmt.Println("Track list is as follows:")
	fmt.Println(GetBBCNowPlaying())
}

var FailCount int = 0

func GetBBCNowPlaying() (out BBCUpdate) {
	r, e := http.Get("http://polling.bbc.co.uk/radio/realtime/bbc_1xtra.jsonp")
	if e != nil {
		FailCount++
		if FailCount > 15 {
			panic("Too many failures from the BBC.")
		}
		return out
	} else {
		text, e := ioutil.ReadAll(r.Body)
		if e != nil {
			FailCount++
			if FailCount > 15 {
				panic("Could not read from the BBC's HTTP stream (wat)")
			}
			return out
		}
		Blank := BBCUpdate{}
		e = json.Unmarshal([]byte(trimBBCoutput(string(text))), &Blank)
		if e != nil {
			FailCount++
			if FailCount > 15 {
				panic("Could not Decode the BBC's json (wat)")
			}
			return out
		}
		return Blank
	}
	return out
}

func trimBBCoutput(input string) string {
	t := strings.Replace(input, "realtimeCallback(", "", 1) // Remove the JSONP thing from the beggining
	return t[:len(t)-1]                                     // And trim the ) off the end.
}

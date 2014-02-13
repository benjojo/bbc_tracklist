package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
	b := flag.String("chan", "bbc_1xtra", "What channel to keep track of")
	flag.Parse()

	starttime := time.Now()
	fmt.Println("Track list is as follows:")
	CurrentTrack := ""
	for {
		Current, e := GetBBCNowPlaying(*b)
		if e == nil {
			if CurrentTrack != fmt.Sprintf("%s - %s", Current.Realtime.Artist, Current.Realtime.Title) {
				CurrentTrack = fmt.Sprintf("%s - %s", Current.Realtime.Artist, Current.Realtime.Title)
				fmt.Printf("%d:%d - %s\n", int(time.Since(starttime).Hours()), int(time.Since(starttime).Minutes())%60, CurrentTrack)
			}

			time.Sleep(time.Second * time.Duration(Current.RequestMinSeconds))
		} else {
			time.Sleep(time.Second * 35)
		}
	}
}

var FailCount int = 0

func GetBBCNowPlaying(url string) (out BBCUpdate, e error) {
	r, e := http.Get(fmt.Sprintf("http://polling.bbc.co.uk/radio/realtime/%s.jsonp", url))
	Failed := fmt.Errorf("Failed to fetch data")
	if e != nil {
		FailCount++
		if FailCount > 15 {
			panic("Too many failures from the BBC.")
		}
		return out, Failed
	} else {
		text, e := ioutil.ReadAll(r.Body)
		if e != nil {
			FailCount++
			if FailCount > 15 {
				panic("Could not read from the BBC's HTTP stream (wat)")
			}
			return out, Failed
		}
		Blank := BBCUpdate{}
		e = json.Unmarshal([]byte(trimBBCoutput(string(text))), &Blank)
		if e != nil {
			FailCount++
			if FailCount > 15 {
				panic("Could not Decode the BBC's json (wat)")
			}
			return out, Failed
		}
		return Blank, e
	}
	return out, Failed
}

func trimBBCoutput(input string) string {
	t := strings.Replace(input, "realtimeCallback(", "", 1) // Remove the JSONP thing from the beggining
	return t[:len(t)-1]                                     // And trim the ) off the end.
}

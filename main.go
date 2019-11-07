package main

import (
	"encoding/xml"
	"fmt"
	"github.com/favish/vmap"
	"log"
	"net/http"
	"os"
	"strconv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	port := getEnv("VMAP_PORT", "80")
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":" + port, nil)
}

var (
	f bool = false
	t bool = true
)

func HelloServer(w http.ResponseWriter, r *http.Request) {

	// Tell the browsers what to do with it
    //TODO: Consider reasonable cache tags, to reduce number of requests from clients

	// Get the required description query param or tell them all is lost.
	dkeys, dok := r.URL.Query()["description_url"]
	if !dok || len(dkeys[0]) < 1 {
		log.Println("")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Url Param 'description_url' is missing"))
		return
	}
	descriptionUrl := dkeys[0]

	// Require referrer to set CORS.
	fkeys, _ := r.URL.Query()["referrer"]
	if len(fkeys) > 0  {
		w.Header().Set("Access-Control-Allow-Origin", fkeys[0])
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Require video duration in order to determine VMAP structure.
	durkeys, dok := r.URL.Query()["duration"]
	if !dok || len(durkeys[0]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Url Param 'duration' is missing"))
		return
	}
    adGapSeconds := 15
	duration, _ := strconv.Atoi(durkeys[0])
	numberOfPods := duration/adGapSeconds
	var adBreaks []vmap.AdBreak

	fmt.Println(w, "Number of pods: %s", numberOfPods)

	if numberOfPods > 0 {
		for i := 1; i <= numberOfPods; i++ {
			adBreaks = append(adBreaks, adBreakGenerator(i * adGapSeconds, descriptionUrl, "midroll", "15", "90", "3"))
		}
	}

	var mainVmap vmap.VMAP = vmap.VMAP{
		Version:    "1.0",
		XmlNS:    "http://www.iab.net/videosuite/vmap",
		AdBreaks:   []vmap.AdBreak{
			{
				TimeOffset:     vmap.Offset{
					Duration: nil,
					Position: vmap.OffsetStart,
					Percent:  0,
				},
				BreakType:      "linear",
				BreakID:        "preroll",
				RepeatAfter:    0,
				AdSource:       &vmap.AdSource{
					ID:               "preroll-ad-1",
					AllowMultipleAds: &f,
					FollowRedirects:  &t,
					VASTAdData:       nil,
					AdTagURI:         &vmap.AdTagURI{
						TemplateType: "vast3",
						URI:          "https://pubads.g.doubleclick.net/gampad/ads?iu=/21841313772/real_vision/preroll&env=vp&impl=s&correlator=&tfcd=0&npa=0&gdfp_req=1&output=vast&sz=640x480&unviewed_position_start=1&description_url=" + descriptionUrl,
					},
					CustomAdData:     nil,
				},
				TrackingEvents: nil,
				Extensions:     nil,
			},
			{
				TimeOffset:     vmap.Offset{
					Duration: nil,
					Position: vmap.OffsetEnd,
					Percent:  0,
				},
				BreakType:      "linear",
				BreakID:        "postroll",
				RepeatAfter:    0,
				AdSource:       &vmap.AdSource{
					ID:               "postroll-ad-1",
					AllowMultipleAds: &f,
					FollowRedirects:  &t,
					VASTAdData:       nil,
					AdTagURI:         &vmap.AdTagURI{
						TemplateType: "vast3",
						URI:          "https://pubads.g.doubleclick.net/gampad/ads?iu=/21841313772/real_vision/postroll&env=vp&impl=s&correlator=&tfcd=0&npa=0&gdfp_req=1&output=vast&sz=640x480&unviewed_position_start=1&description_url="+ descriptionUrl,
					},
					CustomAdData:     nil,
				},
				TrackingEvents: nil,
				Extensions:     nil,
			},
		},
		Extensions: nil,
	}
	mainVmap.AdBreaks = append(mainVmap.AdBreaks, adBreaks...)
	var xmlt, _ = xml.Marshal(mainVmap)
	w.Header().Set("Content-type", "text/xml")
	fmt.Fprintf(w, "%s", xmlt)
}

func adBreakGenerator(offset int, descriptionUrl string, breakId string, minSec string, maxSec string, maxPods string) vmap.AdBreak {
	return vmap.AdBreak {
		TimeOffset:     vmap.Offset{
			Duration: nil,
			Position: offset,
			Percent:  0,
		},
		BreakType:      "linear",
		BreakID:        breakId,
		RepeatAfter:    0,
		AdSource:       &vmap.AdSource{
			ID:                "midroll-ad",
			AllowMultipleAds: &t,
			FollowRedirects:  &t,
			VASTAdData:       nil,
			AdTagURI:         &vmap.AdTagURI{
				TemplateType: "vast3",
				URI:          fmt.Sprintf("https://pubads.g.doubleclick.net/gampad/ads?iu=/21841313772/real_vision/midroll&env=vp&impl=s&correlator=&tfcd=0&npa=0&gdfp_req=1&output=vast&sz=640x480&unviewed_position_start=1&description_url=%s&pmnd=%s&pmxd=%s&pmad=%s",
					descriptionUrl,
					minSec,
					maxSec,
					maxPods,
				),
			},
			CustomAdData:     nil,
		},
		TrackingEvents: nil,
		Extensions:     nil,
	}
}


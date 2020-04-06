package main

import (
	"encoding/xml"
	"fmt"
	"github.com/favish/vmap"
	"github.com/rs/vast"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"
	"net/url"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	port := getEnv("VMAP_PORT", "80")
	http.HandleFunc("/", createVmap)
	http.ListenAndServe(":" + port, nil)
}

var (
	f bool = false
	t bool = true
)

func createVmap(w http.ResponseWriter, r *http.Request) {
    //TODO: Consider reasonable cache tags, to reduce number of requests from clients

	// Referrer to set CORS, explicitly to Google API.
	w.Header().Set("Access-Control-Allow-Origin", "https://imasdk.googleapis.com")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Get the required description query param or error.
	dkeys, dok := r.URL.Query()["description_url"]
	if !dok || len(dkeys[0]) < 1 {
		log.Println("")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Url Param 'description_url' is missing"))
		return
	}
	descriptionUrl := url.QueryEscape(dkeys[0])

	// Get the required referrer if present. It will not always be set.
	rkeys, _ := r.URL.Query()["referrer"]
	referrer := rkeys[0]

	// Grab the partner units based on the referrer passed in. We use different ad units to tell where the traffic is coming from.
    partnerUnitCodes := getPartnerUnit(referrer)

	// Require video duration in order to determine VMAP structure. If this is NaN it means the server does not know the duration and we fallback to a preset VMAP.
	durkeys, dok := r.URL.Query()["duration"]
	if !dok || len(durkeys[0]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Url Param 'duration' is missing"))
		return
	}

    adGapSeconds := 600.0 // An ad pod every 8 minutes.
	durationParameter := durkeys[0]
	duration := 0.0

	// NaN means the video server does not have access to the duration. Assume a 60 minute video so we can cover most use cases.
	if(durationParameter == "NaN") {
		duration = 3600
	} else {
		duration, _ =  strconv.ParseFloat(durationParameter, 32)
	}

	var noAdBeyondPoint = duration - adGapSeconds;
	numberOfPods := duration/adGapSeconds
	var adBreaks []vmap.AdBreak

	if numberOfPods > 1  {
		for i := 1.0; i <= numberOfPods; i++ {
			sec := fmt.Sprintf("%vs", i * adGapSeconds)
			var ter, _ = time.ParseDuration(sec)

			if (i * adGapSeconds) <= noAdBeyondPoint {
				adBreaks = append(adBreaks, adBreakGenerator(vast.Duration(ter), descriptionUrl, "midroll", 0, 30, "3", partnerUnitCodes, i, "midroll"))
			}
		}
	} else if(duration > 300){
		sec := fmt.Sprintf("%vs", 180)
		var ter, _ = time.ParseDuration(sec)
		adBreaks = append(adBreaks, adBreakGenerator(vast.Duration(ter), descriptionUrl, "midroll", 0, 30, "3", partnerUnitCodes, 1, "midroll"))
	}

	// This sets the pre and post roll.
	var mainVmap vmap.VMAP = vmap.VMAP{
		Version:    "1.0",
		XmlNS:    "http://www.iab.net/videosuite/vmap",
		AdBreaks: []vmap.AdBreak{
			{
				TimeOffset: vmap.Offset{
					Duration: nil,
					Position: vmap.OffsetStart,
					Percent:  0,
				},
				BreakType:   "linear",
				BreakID:     "preroll",
				RepeatAfter: 0,
				AdSource: &vmap.AdSource{
					ID:               "preroll-ad-1",
					AllowMultipleAds: &f,
					FollowRedirects:  &t,
					VASTAdData:       nil,
					AdTagURI: &vmap.AdTagURI{
						TemplateType: "vast3",
						URI:          "https://pubads.g.doubleclick.net/gampad/ads?iu=/21841313772/" + partnerUnitCodes[0] + "&env=vp&impl=s&correlator=&tfcd=0&npa=0&gdfp_req=1&output=vast&sz=640x480&unviewed_position_start=1&max_ad_duration=15000&description_url=" + descriptionUrl + "&vpos=preroll&cust_params=testing%3Dtrue&vad_type=linear&pp=real_vision",
					},
					CustomAdData: nil,
				},
				TrackingEvents: nil,
				Extensions:     nil,
			},
			{
				TimeOffset: vmap.Offset{
					Duration: nil,
					Position: vmap.OffsetEnd,
					Percent:  0,
				},
				BreakType:   "linear",
				BreakID:     "postroll",
				RepeatAfter: 0,
				AdSource: &vmap.AdSource{
					ID:               "postroll-ad-1",
					AllowMultipleAds: &f,
					FollowRedirects:  &t,
					VASTAdData:       nil,
					AdTagURI: &vmap.AdTagURI{
						TemplateType: "vast3",
						URI:          "https://pubads.g.doubleclick.net/gampad/ads?iu=/21841313772/" + partnerUnitCodes[2] + "&env=vp&impl=s&correlator=&tfcd=0&npa=0&gdfp_req=1&output=vast&sz=640x480&unviewed_position_start=1&max_ad_duration=15000&description_url=" + descriptionUrl + "&vpos=postroll&cust_params=testing%3Dtrue&vad_type=linear",
					},
					CustomAdData: nil,
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

func getPartnerUnit(referrer string) [3]string {
	var unitCodes [3]string;
	if strings.Contains(referrer, "www.realvision.com") {
		unitCodes[0] = "real_vision/preroll";
		unitCodes[1] = "real_vision/midroll";
		unitCodes[2] = "real_vision/postroll";
	} else if strings.Contains(referrer, "www.zerohedge.com") {
		unitCodes[0] = "real_vision/rv_zerohedge_preroll";
		unitCodes[1] = "real_vision/rv_zerohedge_midroll";
		unitCodes[2] = "real_vision/rv_zerohedge_postroll";
	} else {
		unitCodes[0] = "real_vision/preroll";
		unitCodes[1] = "real_vision/midroll";
		unitCodes[2] = "real_vision/postroll";
	}

	return unitCodes;
}

// This generates optimized ad pods based on the duration of the video.
func adBreakGenerator(offset vast.Duration, descriptionUrl string, breakId string, minSec int, maxSec int, maxPods string, adUnits[3]string, pod float64, vpos string) vmap.AdBreak {
	minSeconds := minSec * 1000
	maxSeconds := maxSec * 1000

	return vmap.AdBreak{
		TimeOffset: vmap.Offset{
			Duration: &offset,
			Position: 0,
			Percent:  0,
		},
		BreakType:   "linear",
		BreakID:     breakId,
		RepeatAfter: 0,
		AdSource: &vmap.AdSource{
			ID:               "midroll-ad",
			AllowMultipleAds: &t,
			FollowRedirects:  &t,
			VASTAdData:       nil,
			AdTagURI: &vmap.AdTagURI{
				TemplateType: "vast3",
				URI: fmt.Sprintf("https://pubads.g.doubleclick.net/gampad/ads?iu=/21841313772/" + adUnits[1] + "&vad_type=linear&env=vp&impl=s&correlator=&tfcd=0&npa=0&gdfp_req=1&output=vast&sz=640x480&unviewed_position_start=1&description_url=%s&pmnd=%v&pmxd=%v&pmad=%v&pod=%v&vpos=%v",
					descriptionUrl,
					minSeconds,
					maxSeconds,
					maxPods,
					pod,
					vpos,
				),
			},
			CustomAdData: nil,
		},
		TrackingEvents: nil,
		Extensions:     nil,
	}
}
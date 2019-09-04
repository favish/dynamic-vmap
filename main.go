package main

import (
	"encoding/xml"
	"fmt"
	"github.com/favish/vmap"
	"log"
	"net/http"
	//"net/url"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":80", nil)
}

var (
	f bool = false
	t bool = true
)

func HelloServer(w http.ResponseWriter, r *http.Request) {

	// Tell the browsers what to do with it
    //TODO: Consider reasonable cache tags, to reduce number of requests from clients

	// Get the required descruption query param or tell them all is lost.
	dkeys, dok := r.URL.Query()["description_url"]
	if !dok || len(dkeys[0]) < 1 {
		log.Println("")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Url Param 'description_url' is missing"))
		return
	}
	descriptionUrl := dkeys[0]

	// also require referrer to set CORS
	fkeys, _ := r.URL.Query()["referrer"]
	if len(fkeys) > 0  {
		//u, _ := url.Parse(fkeys[0])
		w.Header().Set("Access-Control-Allow-Origin", fkeys[0])
	}

	var test vmap.VMAP = vmap.VMAP{
		Version:    "1",
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
						URI:          "https://pubads.g.doubleclick.net/gampad/live/ads?iu=/21778456762/Instream&env=vp&impl=s&correlator=&tfcd=0&npa=0&gdfp_req=1&output=vast&sz=640x480&unviewed_position_start=1&description_url=" + descriptionUrl,
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
						URI:          "https://pubads.g.doubleclick.net/gampad/live/ads?iu=/21778456762/Instream&env=vp&impl=s&correlator=&tfcd=0&npa=0&gdfp_req=1&output=vast&sz=640x480&unviewed_position_start=1&description_url="+ descriptionUrl,
					},
					CustomAdData:     nil,
				},
				TrackingEvents: nil,
				Extensions:     nil,
			},
		},
		Extensions: nil,
	}
	var xmlt, _ = xml.Marshal(test)
	w.Header().Set("Content-type", "text/xml")
	fmt.Fprintf(w, "%s", xmlt)
}

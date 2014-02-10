package plugins

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// For Google's image search API
type googleResult struct{ URL string }
type googleReponseData struct{ Results []googleResult }
type googleResponse struct{ ResponseData googleReponseData }

// Gets a random image URL based on a query string
var Images = Plugin{
	Names: []string{"i", "img", "image", "a", "anim", "animate"},
	Help:  "Returns an image based on a query",
	Run: func(name *string, args *string) (text string) {
		v := url.Values{}
		v.Add("v", "1.0")
		v.Add("q", url.QueryEscape(*args))
		v.Add("safe", "off")

		if *name == "a" || *name == "anim" || *name == "animate" {
			v.Add("as_filetype", "gif")
		}

		uri := "https://ajax.googleapis.com/ajax/services/search/images?" + v.Encode()

		res, err := http.Get(uri)
		if err != nil {
			return "Error finding image"
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "Error reading results"
		}

		var data googleResponse
		if err = json.Unmarshal(body, &data); err != nil {
			return "Error reading JSON"
		}

		return data.ResponseData.Results[0].URL
	},
}

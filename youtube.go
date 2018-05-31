package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func YTInfo(id string) (Video, error) {
	resp, err := http.Get("http://www.youtube.com/get_video_info?&video_id=" + id)

	// fetch the meta information from http
	if err != nil {
		return Video{}, err
	}
	defer resp.Body.Close()

	query_string, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Video{}, err
	}

	return parseMeta(id, string(query_string))
}

func parseMeta(video_id, query_string string) (Video, error) {
	// parse the query string
	u, _ := url.Parse("?" + query_string)

	// parse url params
	query := u.Query()

	// no such video
	if query.Get("errorcode") != "" || query.Get("status") == "fail" {
		return Video{}, errors.New(query.Get("reason"))
	}

	if query.Get("allow_embed") != "1" {
		return Video{}, fmt.Errorf("video cannot be embeded")
	}

	// collate the necessary params
	video := Video{
		ID:    video_id,
		Title: query.Get("title"),
		Img:   query.Get("thumbnail_url"),
	}

	l, _ := strconv.Atoi(query.Get("length_seconds"))
	video.Duration = time.Duration(l) * time.Second

	return video, nil
}

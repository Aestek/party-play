package main

import (
	"sort"
	"sync"
	"time"

	"github.com/knadh/go-get-youtube/youtube"
	"github.com/pkg/errors"
)

var P = &Playlist{
	Index: -1,
	C:     make(chan struct{}),
}

var changedToken = struct{}{}

type Video struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Img      string        `json:"img"`
	Duration time.Duration `json:"duration"`
}

type Item struct {
	lock    sync.RWMutex
	Video   Video   `json:"video"`
	AddedBy *User   `json:"added_by"`
	Likes   []*User `json:"likes"`
}

func (i *Item) Like(by *User) {
	i.lock.Lock()
	defer i.lock.Unlock()

	for _, u := range i.Likes {
		if u.Name == by.Name {
			return
		}
	}

	i.Likes = append(i.Likes, by)
}

type Playlist struct {
	lock             sync.RWMutex
	Items            []*Item       `json:"items"`
	Index            int           `json:"index"`
	CurrentStartedAt time.Time     `json:"current_started_at"`
	C                chan struct{} `json:"-"`
}

func (p *Playlist) Add(id string, by *User) error {
	video, err := youtube.Get(id)
	if err != nil {
		return errors.Wrap(err, "youtube fetch")
	}

	v := Video{
		ID:       id,
		Title:    video.Title,
		Img:      video.Thumbnail_url,
		Duration: time.Duration(video.Length_seconds) * time.Second,
	}

	item := &Item{
		Video:   v,
		AddedBy: by,
		Likes:   []*User{by},
	}

	idx := p.IndexOf(v.ID)
	if idx != -1 {
		p.Like(v.ID, by)
		p.changed()
		return nil
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	p.Items = append(p.Items, item)

	if p.Index != -1 {
		return nil
	}

	p.Index = len(p.Items) - 1
	p.CurrentStartedAt = time.Now()
	p.changed()

	time.AfterFunc(v.Duration, func() {
		p.lock.Lock()
		defer p.lock.Unlock()

		next := p.Index + 1
		if next >= len(p.Items) {
			p.Index = -1
			return
		}

		p.Index = next
		p.CurrentStartedAt = time.Now()
		p.changed()
	})

	return nil
}

func (p *Playlist) start(idx int) {
	if idx >= len(p.Items) {
		p.Index = -1
		return
	}

	p.Index = idx
	p.CurrentStartedAt = time.Now()
	p.changed()

	time.AfterFunc(p.Items[idx].Video.Duration, func() {
		p.lock.Lock()
		defer p.lock.Unlock()

		next := p.Index + 1

		p.Index = next
		p.CurrentStartedAt = time.Now()
		p.changed()
	})
}

func (p *Playlist) Like(id string, by *User) {
	p.lock.Lock()
	defer p.lock.Unlock()

	idx := p.IndexOf(id)
	if idx == -1 {
		return
	}

	p.Items[idx].Like(by)

	sort.SliceStable(p.Items, func(i, j int) bool {
		return len(p.Items[j].Likes) < len(p.Items[i].Likes)
	})
	p.changed()
}

func (p *Playlist) IndexOf(id string) int {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for idx, i := range p.Items {
		if i.Video.ID == id {
			return idx
		}
	}

	return -1
}

func (p *Playlist) changed() {
	p.C <- struct{}{}
}

package main

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/otium/ytdl"

	"github.com/pkg/errors"
)

var P = &Playlist{
	C: make(chan struct{}),
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
	CurrentStartedAt time.Time     `json:"current_started_at"`
	C                chan struct{} `json:"-"`
}

func (p *Playlist) Add(id string, by *User) error {
	vid, err := ytdl.GetVideoInfoFromID(id)
	if err != nil {
		return errors.Wrap(err, "youtube fetch")
	}
	defer p.changed()

	v := Video{
		ID:       id,
		Title:    vid.Title,
		Duration: vid.Duration,
	}

	if v.Duration > 10*time.Minute {
		return fmt.Errorf("Durée de video max: 10min")
	}

	item := &Item{
		Video:   v,
		AddedBy: by,
		Likes:   []*User{by},
	}

	idx := p.IndexOf(v.ID)
	if idx != -1 {
		p.Like(v.ID, by)
		return nil
	}

	for _, i := range p.Items {
		if i.AddedBy.Name == by.Name {
			return fmt.Errorf("Tu as déjà ajouté une vidéo !")
		}
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	l := len(p.Items)
	p.Items = append(p.Items, item)
	if l == 0 {
		p.start()
	}
	return nil
}

func (p *Playlist) start() {
	if len(p.Items) == 0 {
		return
	}

	p.CurrentStartedAt = time.Now()
	time.AfterFunc(p.Items[0].Video.Duration+3*time.Second, p.next)
}

func (p *Playlist) next() {
	p.lock.Lock()
	if len(p.Items) == 0 {
		return
	}
	p.Items = p.Items[1:]
	p.lock.Unlock()
	p.start()
	p.changed()
}

func (p *Playlist) Like(id string, by *User) {
	idx := p.IndexOf(id)
	if idx == -1 {
		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.Items[idx].Like(by)

	if len(p.Items) > 1 {
		sort.SliceStable(p.Items[1:], func(i, j int) bool {
			return len(p.Items[j].Likes) < len(p.Items[i].Likes)
		})
	}
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

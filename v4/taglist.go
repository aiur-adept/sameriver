package sameriver

import (
	"os"
	"sort"

	"encoding/json"
)

type TagList struct {
	tags map[string]bool
	// sorted slice representation is computed lazily
	dirty bool
	slice []string
}

func NewTagList() TagList {
	l := TagList{}
	l.tags = make(map[string]bool)
	return l
}

func (l *TagList) Length() int {
	return len(l.tags)
}

func (l *TagList) Has(tags ...string) bool {
	ok := true
	for _, tag := range tags {
		_, has := l.tags[tag]
		ok = ok && has
	}
	return ok
}

func (l *TagList) Add(tags ...string) {
	if l.tags == nil {
		l.tags = make(map[string]bool, 1)
	}
	for _, t := range tags {
		l.tags[t] = true
	}
	l.dirty = true
}

func (l *TagList) MergeIn(l2 TagList) {
	for t := range l2.tags {
		l.tags[t] = true
	}
	l.dirty = true
}

func (l *TagList) Remove(tag string) {
	delete(l.tags, tag)
	l.dirty = true
}

func (l *TagList) CopyOf() TagList {
	tagsCopy := make(map[string]bool, len(l.tags))
	for tag := range l.tags {
		tagsCopy[tag] = true
	}
	return TagList{
		tags:  tagsCopy,
		dirty: l.dirty,
		slice: l.slice,
	}
}

func (l *TagList) AsSlice() []string {
	if !l.dirty && (l.slice != nil) {
		return l.slice
	} else {
		slice := make([]string, 0, len(l.tags))
		for tag := range l.tags {
			slice = append(slice, tag)
		}
		sort.Strings(slice)
		l.slice = slice
		l.dirty = false
		return l.slice
	}
}

func (l *TagList) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.AsSlice())
}

func (l *TagList) Save(filename string) {
	bytes, err := l.MarshalJSON()
	if err != nil {
		panic(err)
	}
	os.WriteFile(filename, bytes, 0644)
}

func TagListFromJSON(obj []interface{}) TagList {
	l := NewTagList()
	for _, tag := range obj {
		l.Add(tag.(string))
	}
	return l
}

func TagListFromFile(filename string) TagList {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var obj []interface{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}
	return TagListFromJSON(obj)
}

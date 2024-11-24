package history

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var (
	EMPTY_HISTORY = errors.New("History is empty")
)

type PathType int

const (
	TYPE_PROFILE PathType = iota
	TYPE_LIST
	TYPE_THREAD
	TYPE_CAST
)

type Path struct {
	Type      PathType
	Fid       uint64
	Hash      []byte
	Cursor    int
	ViewStart int
	ViewEnd   int
}

func (p Path) String() string {
	return fmt.Sprintf("%d / %d/%s [%d,%d,%d]", p.Type, p.Fid, "0x"+hex.EncodeToString(p.Hash), p.ViewStart, p.Cursor, p.ViewEnd)
}

type History struct {
	paths  []Path
	maxLen int
}

func New(maxLen int) *History {
	return &History{maxLen: maxLen}
}

func (h *History) String() string {
	var buf strings.Builder
	for _, p := range h.paths {
		buf.WriteString(
			fmt.Sprintf("%d FID=%d HASH=%s VIEW=[%d,%d,%d]\n", p.Type, p.Fid, "0x"+hex.EncodeToString(p.Hash), p.ViewStart, p.Cursor, p.ViewEnd),
		)
	}
	return buf.String()
}

func (h *History) Len() int {
	return len(h.paths)
}

func (h *History) MaxLen() int {
	return h.maxLen
}

func (h *History) Free(slots int) {
	if slots >= len(h.paths) {
		h.paths = nil
	} else {
		h.paths = h.paths[slots:]
	}
}

func (h *History) Push(path Path) {
	if len(h.paths) >= h.maxLen {
		h.paths = h.paths[1:]
	}
	h.paths = append(h.paths, path)
}

func (h *History) Pop() (Path, error) {
	if len(h.paths) == 0 {
		return Path{}, EMPTY_HISTORY
	}
	if len(h.paths) == 1 {
		return h.paths[0], nil
	}
	index := len(h.paths) - 1
	path := h.paths[index]
	h.paths = h.paths[:index]
	return path, nil
}

func (h *History) SetView(cursor, start, end int) {
	if len(h.paths) > 0 {
		h.paths[len(h.paths)-1].Cursor = cursor
		h.paths[len(h.paths)-1].ViewStart = start
		h.paths[len(h.paths)-1].ViewEnd = end
	}
}

func (h *History) Peek() (Path, error) {
	if len(h.paths) == 0 {
		return Path{}, EMPTY_HISTORY
	}
	return h.paths[len(h.paths)-1], nil
}

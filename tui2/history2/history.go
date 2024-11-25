package history2

import (
	"errors"
	"fmt"
	"strings"

	tui "github.com/vrypan/fargo/tui2"
)

var (
	EMPTY_HISTORY = errors.New("History is empty")
)

type History struct {
	paths  []tui.CastsStatus
	maxLen int
}

func New(maxLen int) *History {
	return &History{maxLen: maxLen}
}

func (h *History) String() string {
	var buf strings.Builder
	for i, p := range h.paths {
		buf.WriteString(fmt.Sprintf("%d --> ", i))
		buf.WriteString(p.String())
		buf.WriteString("\n\n")
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

func (h *History) Push(path tui.CastsStatus) {
	if len(h.paths) >= h.maxLen {
		h.paths = h.paths[1:]
	}
	h.paths = append(h.paths, path)
}

func (h *History) Pop() (tui.CastsStatus, error) {
	if len(h.paths) == 0 {
		return tui.CastsStatus{}, EMPTY_HISTORY
	}
	if len(h.paths) == 1 {
		return h.paths[0], nil
	}
	index := len(h.paths) - 1
	path := h.paths[index]
	h.paths = h.paths[:index]
	return path, nil
}

func (h *History) Update(newStatus tui.CastsStatus) {
	if len(h.paths) > 0 {
		h.paths[len(h.paths)-1] = newStatus
	}
}

func (h *History) Peek() (tui.CastsStatus, error) {
	if len(h.paths) == 0 {
		return tui.CastsStatus{}, EMPTY_HISTORY
	}
	return h.paths[len(h.paths)-1], nil
}

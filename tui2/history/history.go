package history

import "errors"

var (
	EMPTY_HISTORY = errors.New("History is empty")
	NO_HISTORY    = errors.New("No history available")
)

type PathType int

const (
	TYPE_PROFILE PathType = iota
	TYPE_LIST
	TYPE_THREAD
	TYPE_CAST
)

type Path struct {
	Type PathType
	Fid  uint64
	Hash []byte
}

type History struct {
	paths  []Path
	maxLen int
}

func New(maxLen int) *History {
	return &History{maxLen: maxLen}
}

func (h *History) Len() int {
	return len(h.paths)
}
func (h *History) MaxLen() int {
	return h.maxLen
}
func (h *History) Free(slots int) {
	if slots > h.Len() {
		h.paths = []Path{}
	} else {
		h.paths = h.paths[slots:]
	}
}

func (h *History) Push(path Path) {
	if len(h.paths) == h.maxLen {
		h.paths = h.paths[1:]
	}
	h.paths = append(h.paths, path)
}

func (h *History) Pop() (Path, error) {
	if len(h.paths) == 0 {
		return Path{}, EMPTY_HISTORY
	}
	path := h.paths[len(h.paths)-1]
	h.paths = h.paths[:len(h.paths)-1]
	return path, nil
}

func (h *History) Peek() (Path, error) {
	if len(h.paths) == 0 {
		return Path{}, EMPTY_HISTORY
	}
	return h.paths[len(h.paths)-1], nil
}

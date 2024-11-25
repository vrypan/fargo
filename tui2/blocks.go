package tui2

import (
	"strings"

	"github.com/vrypan/fargo/fctools"
)

// CastsModel.blocks hold an array of pre-rendered casts
// These functions populate the array.

type castsBlock struct {
	id      string
	text    string
	height  int
	padding int
}

func (m *CastsModel) renderBlocks(hash *fctools.Hash, padding int) {
	if hash == nil && m.casts.Head != (fctools.Hash{}) {
		hash = &m.casts.Head
	}
	if hash != nil {
		m.handleThreadBlocks(hash, padding)
	} else {
		m.handleListBlocks(padding)
	}
}

func (m *CastsModel) handleThreadBlocks(hash *fctools.Hash, padding int) {
	idx := m.cursor
	m.hashIdx[idx] = *hash
	text := m.fmtCast(idx, 0)
	m.blocks[idx] = castsBlock{
		id:      hash.String(),
		text:    text,
		height:  strings.Count(text, "\n") + 1,
		padding: padding,
	}
	m.cursor++
	for _, reply := range m.casts.Messages[*hash].Replies {
		m.handleThreadBlocks(&reply, padding+4)
	}
}

func (m *CastsModel) handleListBlocks(padding int) {
	for i, hash := range m.casts.Ordered {
		m.hashIdx[i] = hash
		text := m.fmtCast(i, 0)

		m.blocks[i] = castsBlock{
			id:      hash.String(),
			text:    text,
			height:  strings.Count(text, "\n") + 1,
			padding: padding,
		}
	}
}

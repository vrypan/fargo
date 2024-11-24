package tui2

import (
	"fmt"
	"strconv"
)

func (m *CastsModel) SetFocus(onoff bool, idx int) {
	m.focus = onoff
	m.activeField = 0
	m.cursor = idx
}
func (m *CastsModel) IsFocus() bool {
	return m.focus
}

func (m *CastsModel) GetItemInFocus() string {
	castHash := m.hashIdx[m.cursor]
	message := m.casts.Messages[castHash].Message
	itemCount := len(message.GetData().GetCastAddBody().Mentions) + len(message.GetData().GetCastAddBody().Embeds) + 1
	items := make([]string, itemCount+1)
	items[1] = "fid:" + strconv.FormatUint(message.Data.Fid, 10)

	i := 1
	for _, fid := range message.GetData().GetCastAddBody().Mentions {
		i++
		items[i] = "fid:" + strconv.FormatUint(fid, 10)
	}
	for _, embed := range message.GetData().GetCastAddBody().Embeds {
		i++
		if embedData := embed.GetCastId(); embedData != nil {
			items[i] = fmt.Sprintf("cst:%d:%x", embedData.Fid, embedData.Hash)
		} else {
			items[i] = fmt.Sprintf("url:%s", embed.GetUrl())
		}
	}
	return items[m.activeField]
}

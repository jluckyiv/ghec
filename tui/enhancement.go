package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/jluckyiv/ghec"
)

type item struct {
	be ghec.BaseEnhancement
}

func newItem(be ghec.BaseEnhancement) list.Item {
	return item{be}
}

func (i item) Title() string       { return ghec.Title(i.be) }
func (i item) Description() string { return ghec.Description(i.be) }
func (i item) FilterValue() string { return ghec.Title(i.be) + ghec.Description(i.be) }

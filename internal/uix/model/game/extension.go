package model

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/dustin/go-humanize"
	viewmodel "github.com/seternate/go-lanty-client/internal/ui/model"
)

var _ viewmodel.ChangeNotifier = (*ExtensionModel)(nil)

type ExtensionModel struct {
	Filename   string
	Name       binding.String
	Size       binding.String
	Downloaded binding.Bool
}

func NewExtensionModel(filename, name string, size uint64, downloaded bool) *ExtensionModel {
	ext := &ExtensionModel{
		Filename:   filename,
		Name:       binding.NewString(),
		Size:       binding.NewString(),
		Downloaded: binding.NewBool(),
	}
	ext.Name.Set(name)
	if size > 0 {
		ext.Size.Set(humanize.Bytes(size))
	}
	ext.Downloaded.Set(downloaded)
	return ext
}

func (m *ExtensionModel) SetName(name string) {
	m.Name.Set(name)
}

func (m *ExtensionModel) SetDownloaded(downloaded bool) {
	m.Downloaded.Set(downloaded)
}

func (m *ExtensionModel) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)
	m.Name.AddListener(l)
	m.Size.AddListener(l)
	m.Downloaded.AddListener(l)
}

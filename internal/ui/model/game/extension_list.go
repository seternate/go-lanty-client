package model

import (
	"sync"

	viewmodel "github.com/seternate/go-lanty-client/internal/ui/model"
)

var _ viewmodel.ChangeNotifier = (*ExtensionListModel)(nil)

type ExtensionListModel struct {
	mu        sync.RWMutex
	extensions []*ExtensionModel
	listeners []func()
}

func NewExtensionListModel() *ExtensionListModel {
	return &ExtensionListModel{
		extensions: nil,
		listeners:  nil,
	}
}

func (m *ExtensionListModel) AddExtension(ext *ExtensionModel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.extensions = append(m.extensions, ext)
	ext.AddChangeListener(m.notify)
	m.notifyLocked()
}

func (m *ExtensionListModel) GetExtensions() []*ExtensionModel {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*ExtensionModel, len(m.extensions))
	copy(result, m.extensions)
	return result
}

func (m *ExtensionListModel) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.extensions = nil
	m.notifyLocked()
}

func (m *ExtensionListModel) AddChangeListener(fn func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners = append(m.listeners, fn)
}

func (m *ExtensionListModel) notify() {
	m.mu.RLock()
	listeners := make([]func(), len(m.listeners))
	copy(listeners, m.listeners)
	m.mu.RUnlock()
	for _, fn := range listeners {
		fn()
	}
}

func (m *ExtensionListModel) notifyLocked() {
	listeners := make([]func(), len(m.listeners))
	copy(listeners, m.listeners)
	m.mu.Unlock()
	for _, fn := range listeners {
		fn()
	}
	m.mu.Lock()
}

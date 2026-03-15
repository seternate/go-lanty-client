package viewmodel

import "fyne.io/fyne/v2/data/binding"

type HeaderScrollScreen struct {
	title binding.String
	info  binding.String
}

func NewHeaderScrollScreen(title string) *HeaderScrollScreen {
	vm := &HeaderScrollScreen{
		title: binding.NewString(),
		info:  binding.NewString(),
	}

	vm.title.Set(title)

	return vm
}

func NewHeaderScrollScreenWithData(title string, info binding.String) *HeaderScrollScreen {
	vm := &HeaderScrollScreen{
		title: binding.NewString(),
		info:  info,
	}

	vm.title.Set(title)

	return vm
}

func (vm *HeaderScrollScreen) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)

	vm.title.AddListener(l)
	vm.info.AddListener(l)
}

func (vm *HeaderScrollScreen) GetTitle() string {
	title, err := vm.title.Get()
	if err != nil {
		return ""
	}
	return title
}

func (vm *HeaderScrollScreen) GetInfo() string {
	info, err := vm.info.Get()
	if err != nil {
		return ""
	}
	return info
}

func (vm *HeaderScrollScreen) SetInfo(info string) {
	vm.info.Set(info)
}

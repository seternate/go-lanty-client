package viewmodel

type ChangeNotifier interface {
	AddChangeListener(fn func())
}

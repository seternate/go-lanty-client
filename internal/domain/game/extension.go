package game

type Extension struct {
	GameSlug   string
	Filename   string
	Name       string
	Size       uint64
	Downloaded bool
}

func NewExtension(gameSlug, filename, name string, size uint64, downloaded bool) Extension {
	return Extension{
		GameSlug:   gameSlug,
		Filename:   filename,
		Name:       name,
		Size:       size,
		Downloaded: downloaded,
	}
}

func (e *Extension) MarkDownloaded() {
	e.Downloaded = true
}

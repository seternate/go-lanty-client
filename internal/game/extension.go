package game

type Extension struct {
	Slug       string
	Filename   string
	Name       string
	Size       uint64
	Downloaded bool
}

func NewExtension(slug, filename, name string, size uint64, downloaded bool) Extension {
	return Extension{
		Slug:       slug,
		Filename:   filename,
		Name:       name,
		Size:       size,
		Downloaded: downloaded,
	}
}

func (e *Extension) MarkDownloaded() {
	e.Downloaded = true
}

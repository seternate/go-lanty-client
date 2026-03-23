package game

type LaunchParamGroup struct {
	CategoryKey string
	Label       string
	Params      []LaunchParam
}

func (group LaunchParamGroup) HeaderText() string {
	if group.CategoryKey == "" {
		return "General"
	}
	return group.Label
}

func GroupLaunchParamsByCategory(params []LaunchParam) []LaunchParamGroup {
	var sections []LaunchParamGroup
	keyToIndex := make(map[string]int)
	for _, p := range params {
		key := p.Category()
		label := key
		if key == "" {
			label = ""
		}
		idx, ok := keyToIndex[key]
		if !ok {
			keyToIndex[key] = len(sections)
			sections = append(sections, LaunchParamGroup{
				CategoryKey: key,
				Label:       label,
			})
			idx = len(sections) - 1
		}
		sec := &sections[idx]
		sec.Params = append(sec.Params, p)
	}
	return sections
}

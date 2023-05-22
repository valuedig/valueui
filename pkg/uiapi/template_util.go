package uiapi

func (it *TemplateLayout) Refix() *TemplateLayout {
	if it.Width == "" {
		it.Width = "auto"
	}
	if it.Options == nil {
		it.Options = map[string]string{}
	}
	if len(it.Cols) == 0 {
		it.Cols = append(it.Cols, &TemplateLayout{
			Width: "auto",
		})
	}
	return it
}

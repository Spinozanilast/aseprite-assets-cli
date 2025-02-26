package cmddescription

type Builder struct {
	styles      *DescriptionStyles
	description *CmdDescription
}

func NewBuilder() *Builder {
	return &Builder{
		styles:      DefaultStyles(),
		description: &CmdDescription{},
	}
}

func (b *Builder) WithStyles(styles *DescriptionStyles) *Builder {
	b.styles = styles
	return b
}

func (b *Builder) WithTitle(title string) *Builder {
	b.description.Title = title
	return b
}

func (b *Builder) WithShort(short string) *Builder {
	b.description.Short = short
	return b
}

func (b *Builder) WithLong(long string) *Builder {
	b.description.Long += "\n" + long
	return b
}

func (b *Builder) WithSection(title string, text string) *Builder {
	b.description.Sections = append(b.description.Sections, CmdDescriptionSection{
		Title: title,
		Text:  text,
	})
	return b
}

// WithList adds a list to the long description taking items from items string args
func (b *Builder) WithList(title string, indent int, indentChar string, items ...string) *Builder {
	list := CmdDescriptionList{
		Title:      title,
		Items:      make([]CmdDescriptionListItem, len(items)),
		Indent:     indent,
		IndentChar: indentChar,
	}

	for i, item := range items {
		list.Items[i] = CmdDescriptionListItem{ItemName: item}
	}

	list.Construct(b.styles)
	b.description.Long += list.String()
	return b
}

// WithList adds a list to the long description taking items from CmdDescriptionListItem args
func (b *Builder) WithListFromItems(title string, indent int, indentChar string, items ...*CmdDescriptionListItem) *Builder {
	list := CmdDescriptionList{
		Title:      title,
		Items:      make([]CmdDescriptionListItem, len(items)),
		Indent:     indent,
		IndentChar: indentChar,
	}

	for i, item := range items {
		list.Items[i] = *item
	}

	list.Construct(b.styles)
	b.description.Long += list.String()
	return b
}

// Build constructs and returns the final command description
func (b *Builder) Build() *CmdDescription {
	b.description.Construct(b.styles)
	return b.description
}

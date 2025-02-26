package cmddescription

import (
	"strings"
)

type CmdDescription struct {
	Title          string
	Short          string
	Long           string
	Sections       []CmdDescriptionSection
	SectionDivider string
}

type CmdDescriptionSection struct {
	Title string
	Text  string
}

type CmdDescriptionList struct {
	Title      string
	Indent     int
	IndentChar string
	Items      []CmdDescriptionListItem
}

type CmdDescriptionListItem struct {
	ItemName    string
	Description string
	SubItems    []CmdDescriptionListItem
}

func (c *CmdDescription) Construct(styles *DescriptionStyles) {
	c.Title = styles.Title.Render(c.Title)
	for i := range c.Sections {

		sectionTitle := styles.SectionTitle.Render(c.Sections[i].Title)
		sectionText := styles.SectionText.Render(c.Sections[i].Text)

		c.Sections[i].Title = sectionTitle
		c.Sections[i].Text = sectionText
	}
}

func (l *CmdDescriptionList) Construct(styles *DescriptionStyles) {
	l.Title = styles.ListTitle.Render(l.Title)
	for i := range l.Items {
		l.Items[i] = l.Items[i].constructHierarchy(styles, 0, l.Indent, l.IndentChar)
	}
}

func (li *CmdDescriptionListItem) constructHierarchy(styles *DescriptionStyles, depth int, indents int, indentChar string) CmdDescriptionListItem {
	indent := strings.Repeat(" ", depth*indents)
	descriptionIndent := indent + strings.Repeat(" ", len(indentChar))

	var styledName string
	var styledDescription string
	if depth == 0 {
		styledName = styles.ListItem.Render(li.ItemName)
		if li.Description != "" {
			styledDescription = styles.ListItemDescription.Render("# " + li.Description)
		}
	} else {
		styledName = styles.ListSubItem.Render(li.ItemName)
		if li.Description != "" {
			styledDescription = styles.ListItemDescription.Render("# " + li.Description)
		}
	}

	li.ItemName = indent + indentChar + styledName
	if li.Description != "" {
		li.Description = styledDescription + descriptionIndent
	}

	for i := range li.SubItems {
		li.SubItems[i] = li.SubItems[i].constructHierarchy(styles, depth+1, indents, indentChar)
	}

	return *li
}

func (c *CmdDescription) String() string {
	var sb strings.Builder

	sb.WriteString(c.Title + "\n")

	if c.Long != "" {
		sb.WriteString(c.Long + "\n\n")
	}

	if len(c.Sections) > 0 {
		for i, section := range c.Sections {
			sb.WriteString(section.String())
			if i < len(c.Sections)-1 && c.SectionDivider != "" {
				sb.WriteString(c.SectionDivider + "\n")
			}
		}
	}

	return sb.String()
}

func (s *CmdDescriptionSection) String() string {
	var sb strings.Builder

	sb.WriteString(s.Title + "\n")
	sb.WriteString(s.Text + "\n\n")

	return sb.String()
}

func (l *CmdDescriptionList) String() string {
	var sb strings.Builder

	sb.WriteString(l.Title + "\n")

	for _, item := range l.Items {
		sb.WriteString(item.String())
	}

	return sb.String()
}

func (li *CmdDescriptionListItem) String() string {
	var sb strings.Builder

	if li.Description != "" {
		sb.WriteString(li.Description + "\n" + li.ItemName + "\n")
	} else {
		sb.WriteString(li.ItemName + "\n")
	}

	for _, subItem := range li.SubItems {
		sb.WriteString(subItem.String())
	}

	return sb.String()
}

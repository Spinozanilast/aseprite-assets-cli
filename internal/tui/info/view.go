package info

import (
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/spinozanilast/aseprite-assets-cli/pkg/consts"
)

const (
	maxPathLength = 45
	dateFormat    = "Jan 02, 2006 15:04"
	kb            = 1024
	mb            = kb * kb
)

func (m Model) View() string {
	content := m.Styles.Base.Width(m.Width - 4)

	switch {
	case m.Error != "":
		return content.Render(m.Styles.Error.Render(m.Error))
	case m.AssetInfo == nil:
		return content.Render(m.Styles.NoContent.Render("No asset selected"))
	default:
		return content.Render(m.renderAssetInfo())
	}
}

func (m Model) renderAssetInfo() string {
	var sb strings.Builder

	sb.WriteString(m.renderTitle())
	sb.WriteString("\n\n")
	sb.WriteString(m.renderMetadata())
	sb.WriteString(m.renderPath())

	return sb.String()
}

func (m Model) renderTitle() string {
	return m.Styles.Title.Render(m.AssetInfo.Name)
}

func (m Model) renderMetadata() string {
	return strings.Join([]string{
		m.renderField("Type:", string(m.AssetInfo.Type)),
		m.renderField("Size:", formatSize(m.AssetInfo.Size)),
		m.renderField("Extension:", m.AssetInfo.Extension),
		m.renderField("Modified:", m.AssetInfo.ModTime.Format(dateFormat)),
	}, "\n") + "\n"
}

func (m Model) renderField(label, value string) string {
	return fmt.Sprintf("%s %s",
		m.Styles.Label.Render(label),
		m.Styles.Value.Render(value),
	)
}

func (m Model) renderPath() string {
	truncatedPath := truncateMiddle(m.AssetInfo.Path, maxPathLength)
	return m.renderField("Path:", truncatedPath)
}

func formatSize(size int64) string {
	switch {
	case size < kb:
		return fmt.Sprintf("%d B", size)
	case size < mb:
		return fmt.Sprintf("%.2f KB", float64(size)/kb)
	default:
		return fmt.Sprintf("%.2f MB", float64(size)/mb)
	}
}

func truncateMiddle(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}

	sep := "..."
	base := filepath.Base(path)
	remaining := maxLen - len(base) - len(sep)
	if remaining < 1 {
		return sep + base
	}

	dir := filepath.Dir(path)
	head := strings.Split(dir, string(filepath.Separator))
	var builder strings.Builder

	for i := len(head) - 1; i >= 0; i-- {
		if builder.Len()+len(head[i])+len(sep) > remaining {
			break
		}
		if builder.Len() > 0 {
			builder.WriteString(string(filepath.Separator))
		}
		builder.WriteString(head[i])
	}

	truncated := filepath.Join(
		reversePath(builder.String()),
		sep,
		base,
	)

	if len(truncated) > maxLen {
		return sep + base
	}
	return truncated
}

func reversePath(path string) string {
	parts := strings.Split(path, string(filepath.Separator))
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return filepath.Join(parts...)
}

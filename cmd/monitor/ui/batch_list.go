package ui

import (
	"strings"
	"time"
	"sort"
	"github.com/charmbracelet/lipgloss"
	openai "github.com/sashabaranov/go-openai"
)

type batchItem struct {
	id        string
	status    string
	createdAt time.Time
	counts    openai.BatchRequestCounts
}

func processBatches(batches []openai.BatchResponse) []batchItem {
    items := make([]batchItem, len(batches))
    for i, b := range batches {
        items[i] = batchItem{
            id:        b.ID,
            status:    b.Status,
            createdAt: time.Unix(int64(b.CreatedAt), 0),
            counts:    b.RequestCounts,
        }
    }

    // Sort by createdAt in descending order (newest first)
    sort.Slice(items, func(i, j int) bool {
        return items[i].createdAt.After(items[j].createdAt)
    })

    return items
}

func (m Model) filterBatches() []batchItem {
	var filtered []batchItem
	for _, b := range m.batches {
		switch m.currentTab {
		case activeTab:
			if b.status == "in_progress" {
				filtered = append(filtered, b)
			}
		case completedTab:
			if b.status == "completed" {
				filtered = append(filtered, b)
			}
		case failedTab:
			if b.status == "failed" {
				filtered = append(filtered, b)
			}
		case expiredTab:
            if b.status == "expired" {
                filtered = append(filtered, b)
            }
		}
		
	}
	return filtered
}

func (m Model) renderBatches(batches []batchItem) string {
    if len(batches) == 0 {
        return "No batches found"
    }

    // Render header row
    headerRow := lipgloss.JoinHorizontal(
        lipgloss.Left,
        batchIDStyle.Copy().Render("Batch ID"),
        "  ",
        batchIDStyle.Copy().Render("Submitted"),
        "  ",
        batchIDStyle.Copy().Render("Status"),
    )

    separator := lipgloss.NewStyle().
        Foreground(primaryColor).
        Render(strings.Repeat("─", m.width-4))

    var rendered []string
    rendered = append(rendered, headerRow, separator)

    // Render each batch
    for i, batch := range batches {
        absoluteIndex := m.offset + i
        var style lipgloss.Style
        if absoluteIndex == m.cursor {
            style = selectedBatchStyle // Use the new bold style for selected batch
        } else {
            style = lipgloss.NewStyle()
        }

        // Batch info line
        batchInfo := lipgloss.JoinHorizontal(
            lipgloss.Left,
            batch.id,
            "  ",
            batch.createdAt.Format("2006-01-02 15:04:05"),
            "  ",
            statusStyle[batch.status].Render(batch.status),
        )

        progress := renderProgressWithCounts(
            batch.counts.Completed, 
            batch.counts.Total, 
            m.width-4,
        )
        
        batchBlock := lipgloss.JoinVertical(
            lipgloss.Left,
            batchInfo,
            progress,
        )
        
        rendered = append(rendered, style.Render(batchBlock))
    }
    
    return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

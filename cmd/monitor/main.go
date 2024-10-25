package main

import (
	"batch-gpt/cmd/monitor/ui"
	"batch-gpt/server/db"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
    // Initialize MongoDB connection
    db.InitMongoDB()

    p := tea.NewProgram(
        ui.NewModel(),
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )

    if _, err := p.Run(); err != nil {
        fmt.Printf("Error running program: %v", err)
        os.Exit(1)
    }
}

package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
    // Theme colors
    primaryColor    = lipgloss.Color("86")    // cyan
    secondaryColor  = lipgloss.Color("87")    // light cyan
    alertColor      = lipgloss.Color("196")   // red
    warningColor    = lipgloss.Color("214")   // orange
    successColor    = lipgloss.Color("46")    // green
    textColor       = lipgloss.Color("252")   // light gray

    // Styles
    titleStyle = lipgloss.NewStyle().
            Bold(true).
            Foreground(primaryColor).
            BorderStyle(lipgloss.RoundedBorder()).
            BorderForeground(primaryColor).
            Padding(0, 1).
            Align(lipgloss.Center)

    tabStyle = lipgloss.NewStyle().
        Padding(0, 1)

    activeTabStyle = tabStyle.Copy().
            Foreground(primaryColor).
            Bold(true). // Make active tab bold
            Border(lipgloss.RoundedBorder(), false, false, true, false).
            BorderForeground(primaryColor)
    
    selectedBatchStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(primaryColor).
        Bold(true). // Make selected batch bold
        Padding(0, 1)

    batchIDStyle = lipgloss.NewStyle().
        Foreground(secondaryColor)

    statusStyle = map[string]lipgloss.Style{
        "completed": lipgloss.NewStyle().Foreground(successColor),
        "failed":    lipgloss.NewStyle().Foreground(alertColor),
        "in_progress": lipgloss.NewStyle().Foreground(warningColor),
    }

    progressBarStyle = lipgloss.NewStyle().
        Foreground(primaryColor)

    progressBarFilledStyle = lipgloss.NewStyle().
        Background(primaryColor)

    progressBarErrorStyle = lipgloss.NewStyle().
        Background(alertColor)

    helpStyle = lipgloss.NewStyle().
        Foreground(textColor).
        Align(lipgloss.Right)
)

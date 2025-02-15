// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"context"
	"errors"
	"fmt"

	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	views_util "github.com/daytonaio/daytona/pkg/views/util"
)

const maxWidth = 160

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(0, 4, 1, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(views.Green).
		Bold(true).
		Padding(1, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(views.Green).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(views.Green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(views.Green)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type Model struct {
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	width  int
}

func GetRepositoryFromUrlInput(multiWorkspace bool, workspaceOrder int, apiClient *apiclient.APIClient, selectedRepos map[string]int) (*apiclient.GitRepository, error) {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	title := "Git repository"

	if multiWorkspace {
		title = getOrderNumberString(workspaceOrder) + " workspace repository"
	}

	var initialRepoUrl string
	var repo *apiclient.GitRepository

	var previousRepoUrl string
	var previousError error

	initialRepoInput := huh.NewInput().
		Title(title).
		Value(&initialRepoUrl).
		Key("initialWorkspaceRepo").
		Validate(func(str string) error {
			if str == previousRepoUrl && previousError != nil {
				return previousError
			}

			previousRepoUrl = str
			previousError = nil

			err := views_util.WithInlineSpinner("Validating", func() error {
				var err error
				repo, err = validateRepoUrl(str, apiClient)
				return err
			})

			if err != nil {
				previousError = err
			}

			return err
		})

	dTheme := views.GetCustomTheme()

	m.form = huh.NewForm(
		huh.NewGroup(
			initialRepoInput,
		).WithHeight(5),
	).WithTheme(dTheme).
		WithWidth(maxWidth).
		WithShowHelp(false).
		WithShowErrors(true)

	err := m.form.WithProgramOptions(tea.WithAltScreen()).Run()
	if err != nil {
		return nil, err
	}

	selectedRepos[repo.Url]++

	return repo, nil
}

func RunAddMoreWorkspacesForm() (bool, error) {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	var addMore bool

	confirmInput :=
		huh.NewConfirm().
			Title("Add another workspace?").
			Value(&addMore)

	m.form = huh.NewForm(
		huh.NewGroup(confirmInput),
	).
		WithWidth(maxWidth).
		WithShowHelp(false).
		WithShowErrors(true).
		WithTheme(views.GetCustomTheme())

	err := m.form.WithProgramOptions(tea.WithAltScreen()).Run()
	if err != nil {
		return false, err
	}

	return addMore, nil
}

func getOrderNumberString(number int) string {
	if number >= 1 && number <= 10 {
		// Handle numbers 1 to 10
		switch number {
		case 1:
			return "First"
		case 2:
			return "Second"
		case 3:
			return "Third"
		case 4:
			return "Fourth"
		case 5:
			return "Fifth"
		case 6:
			return "Sixth"
		case 7:
			return "Seventh"
		case 8:
			return "Eighth"
		case 9:
			return "Ninth"
		case 10:
			return "Tenth"
		}
	} else if number >= 11 {
		// Handle numbers 11 and beyond
		return fmt.Sprintf("%d.", number)
	}
	// Handle invalid numbers or negative numbers
	return "Invalid"
}

func validateRepoUrl(repoUrl string, apiClient *apiclient.APIClient) (*apiclient.GitRepository, error) {
	result, err := util.GetValidatedUrl(repoUrl)
	if err != nil {
		return nil, err
	}

	repo, _, err := apiClient.GitProviderAPI.GetGitContext(context.Background()).Repository(apiclient.GetRepositoryContext{
		Url: result,
	}).Execute()
	if err != nil {
		wrappedErr := "Failed to fetch repository information. Please check the URL and try again."
		return nil, errors.New(views_util.WrapText(wrappedErr, views_util.GetTerminalWidth()))
	}

	return repo, nil
}

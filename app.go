package main

import (
	"context"
	"fmt"
	"goinload/internal/scraper"
	"goinload/internal/utils"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SendUrl(link string) {
	matchedDomains, isValid := utils.IsValidURL(link)

	if isValid {
		utils.EmitMsg(a.ctx, utils.InfoTask, fmt.Sprintf("Supported domain! Wait a moment %v...", matchedDomains), false)

		utils.EmitMsg(a.ctx, utils.Status, "started", false)

		scraper.Scrape(a.ctx, matchedDomains, link)

		utils.EmitMsg(a.ctx, utils.Status, "stopped", false)
	} else {
		utils.EmitMsg(a.ctx, utils.InfoTask, "Invalid url (accepted ex: https://cyberdrop.me/a/XXXXXXX)", false)
	}
}

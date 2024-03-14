package scraper

import (
	"context"
	"fmt"
	"goinload/internal/scraper/targets"
	"goinload/internal/utils"
)

// Scrape function accepts a link to scrape
func Scrape(ctx context.Context, target utils.TargetDomain, link string) error {

	switch target {
	case utils.Cyberdrop:
		targets.TargetCyberdrop(ctx, link)

	case utils.Ososedki:
		targets.TargetOsosedki(ctx, link)

	case utils.Bunkr:
		targets.TargetBunkr(ctx, link)

	default:
		return fmt.Errorf("pattern match failed for %s: %s", target, link)
	}

	return nil
}

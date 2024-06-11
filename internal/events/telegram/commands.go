package telegram2

import (
	e "adviserbot/internal/lib/err"
	"adviserbot/internal/storage"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text, username string, chatID int) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command %s from %s", text, username)
	if isAddCmd(text) {
		return p.SavePage(chatID, text, username)
	}
	switch text {
	case RndCmd:
		return p.SendRandom(chatID, username)
	case HelpCmd:
		return p.SendHelp(chatID)
	case StartCmd:
		return p.SendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

}

func (p *Processor) SavePage(chatID int, text string, username string) error {
	const op = "events.telegram.SavePage"
	page := storage.Page{
		URL:      text,
		UserName: username,
	}
	ok, err := p.storage.IsExists(&page)
	if err != nil {
		return e.Wrap(op, err)
	}
	if ok {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(&page); err != nil {
		return e.Wrap(op, err)
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return e.Wrap(op, err)
	}
	return nil
}

func (p *Processor) SendRandom(chatID int, username string) error {
	const op = "events.telegram.SendRandom"
	page, err := p.storage.PickRandom(username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return p.tg.SendMessage(chatID, msgNoSavedPages)
		}
		return e.Wrap(op, err)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return e.Wrap(op, err)
	}
	return p.storage.Remove(page)
}

func (p *Processor) SendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) SendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

package telegram2

import (
	"adviserbot/internal/clients/telegram"
	"adviserbot/internal/events"
	e "adviserbot/internal/lib/err"
	"adviserbot/internal/storage"
	"errors"
)

var (
	ErrUnknownType = errors.New("unknown type")
	ErrCantGetMeta = errors.New("can't get meta from event")
)

type Processor struct {
	tg      *telegram.CLient
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

func New(client *telegram.CLient, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	const op = "events.telegram.Fetch"
	updates, err := p.tg.Update(p.offset, limit)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	if len(updates) == 0 {
		return nil, nil
	}
	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	const op = "events.telegram.Process"
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap(op, ErrUnknownType)
	}

}

func (p *Processor) processMessage(event events.Event) error {
	const op = "events.telegram.processMessage"
	meta, err := GetMeta(event)
	if err != nil {
		return e.Wrap(op, err)
	}
	if err := p.doCmd(event.Text, meta.UserName, meta.ChatID); err != nil {
		return e.Wrap(op, err)
	}
	return nil
}

func GetMeta(event events.Event) (Meta, error) {
	const op = "events.telegram.meta"
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap(op, ErrCantGetMeta)
	}
	return res, nil
}

func event(u telegram.Update) events.Event {
	uType := fetchType(u)
	res := events.Event{
		Type: uType,
		Text: fetchText(u),
	}
	if uType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			UserName: u.Message.From.Username,
		}
	}
	return res

}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message

}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}

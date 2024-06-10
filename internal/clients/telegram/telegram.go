package telegram

import (
	e "adviserbot/internal/lib/err"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type CLient struct {
	host     string
	basePath string
	client   http.Client
}

const (
	UpdateMethod      = "getUpdate"
	SendMessageMethod = "sendMessage"
)

func New(host, token string) *CLient {
	return &CLient{
		host:     host,
		basePath: fmt.Sprintf("bot%s", token),
		client:   http.Client{},
	}
}

func (c *CLient) Update(offset, limit int) ([]Update, error) {
	const op = "clients.telegram.Update"
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	data, err := c.doRequest(q, UpdateMethod)
	if err != nil {
		return []Update{}, e.Wrap(op, err)
	}
	var res UpdateResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []Update{}, e.Wrap(op, err)
	}
	return res.Result, nil
}

func (c *CLient) SendMessage(chatId int, text string) error {
	const op = "clients.telegram.SendMessage"
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(q, SendMessageMethod)
	if err != nil {
		return e.Wrap(op, err)
	}
	return nil

}

func (c *CLient) doRequest(query url.Values, method string) ([]byte, error) {
	const op = "clients.telegram.doRequest"
	url := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	return body, nil
}

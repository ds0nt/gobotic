package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

const telegramBaseURL = "https://api.telegram.org"

type telegramClient struct {
	token  string
	client *http.Client
	Debug  bool
}

func newTelegramClient(token string) *telegramClient {

	if token == "" {
		panic("empty token not allowed")
	}
	return &telegramClient{
		token:  token,
		client: http.DefaultClient,
	}
}

func (c *telegramClient) makeURL(url string) string {
	return fmt.Sprintf("%s/bot%s/%s", telegramBaseURL, c.token, url)
}

func (c *telegramClient) doRequest(method, url string, payload interface{}) (resp *http.Response, baseResp *telegramBaseResponse, err error) {
	if c.Debug {
		pretty.Println(payload)
	}

	buf := bytes.Buffer{}
	w := json.NewEncoder(&buf)
	if err = w.Encode(payload); err != nil {
		err = errors.Wrap(err, "error encoding payload")
		return
	}
	req, err := http.NewRequest(method, c.makeURL(url), &buf)
	if err != nil {
		err = errors.Wrap(err, "error during new http request")
		return
	}
	req.Header.Set("content-type", "application/json")
	if c.Debug {
		pretty.Println(req.URL)
	}

	resp, err = c.client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "error performing request")
		return
	}
	defer func() {
		_err := resp.Body.Close()
		if _err != nil {
			fmt.Println(errors.Wrap(err, "error closing response body"))
			debug.PrintStack()
		}
	}()

	d := json.NewDecoder(resp.Body)
	baseResp = &telegramBaseResponse{}
	if err = d.Decode(baseResp); err != nil {
		err = errors.Wrap(err, "error reading response body")
		return
	}
	if !baseResp.OK {
		err = errors.New("telegram response body field OK was false")
		if c.Debug {
			pretty.Println(baseResp)
		}
		return
	}

	return resp, baseResp, nil
}

type telegramBaseResponse struct {
	OK          bool            `json:"ok"`
	Description string          `json:"description"`
	Result      json.RawMessage `json:"result"`
}

func (c *telegramClient) SetWebhook(url string) error {
	payload := map[string]interface{}{
		"url": url,
	}
	_, _, err := c.doRequest(http.MethodPost, "setWebhook", payload)
	if err != nil {
		return err
	}

	return nil
}
func (c *telegramClient) SendMessage(message *SendMessageRequest) error {
	_, _, err := c.doRequest(http.MethodPost, "sendMessage", message)
	if err != nil {
		return err
	}

	return nil
}

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
	// EditedMessage      *Message            `json:"edited_message"`
	// ChannelPost        *Message            `json:"channel_post"`
	// EditedChannelPost  *Message            `json:"edited_channel_post"`
	// InlineQuery        *InlineQuery        `json:"inline_query"`
	// ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	// CallbackQuery      *CallbackQuery      `json:"callback_query"`
	// ShippingQuery      *ShippingQuery      `json:"shipping_query"`
	// PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query"`
}

// Chat contains information about the place a message was sent.
type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`      // optional
	UserName  string `json:"username"`   // optional
	FirstName string `json:"first_name"` // optional
	LastName  string `json:"last_name"`  // optional
	// AllMembersAreAdmins bool   `json:"all_members_are_administrators"` // optional
	// Photo               *ChatPhoto `json:"photo"`
	Description string `json:"description,omitempty"` // optional
	InviteLink  string `json:"invite_link,omitempty"` // optional
}

type SendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// Message is returned by almost every request, and contains data about
// almost anything.
type Message struct {
	MessageID int   `json:"message_id"`
	From      *User `json:"from"` // optional
	Date      int   `json:"date"`
	Chat      *Chat `json:"chat"`
	// ForwardFrom           *User              `json:"forward_from"`            // optional
	// ForwardFromChat       *Chat              `json:"forward_from_chat"`       // optional
	ForwardFromMessageID int      `json:"forward_from_message_id"` // optional
	ForwardDate          int      `json:"forward_date"`            // optional
	ReplyToMessage       *Message `json:"reply_to_message"`        // optional
	EditDate             int      `json:"edit_date"`               // optional
	Text                 string   `json:"text"`                    // optional
	// Entities              *[]MessageEntity   `json:"entities"`                // optional
	// Audio                 *Audio             `json:"audio"`                   // optional
	// Document              *Document          `json:"document"`                // optional
	// Game                  *Game              `json:"game"`                    // optional
	// Photo                 *[]PhotoSize       `json:"photo"`                   // optional
	// Sticker               *Sticker           `json:"sticker"`                 // optional
	// Video                 *Video             `json:"video"`                   // optional
	// VideoNote             *VideoNote         `json:"video_note"`              // optional
	// Voice                 *Voice             `json:"voice"`                   // optional
	// Caption               string             `json:"caption"`                 // optional
	// Contact               *Contact           `json:"contact"`                 // optional
	// Location              *Location          `json:"location"`                // optional
	// Venue                 *Venue             `json:"venue"`                   // optional
	// NewChatMembers        *[]User            `json:"new_chat_members"`        // optional
	// LeftChatMember        *User              `json:"left_chat_member"`        // optional
	// NewChatTitle          string             `json:"new_chat_title"`          // optional
	// NewChatPhoto          *[]PhotoSize       `json:"new_chat_photo"`          // optional
	// DeleteChatPhoto       bool               `json:"delete_chat_photo"`       // optional
	// GroupChatCreated      bool               `json:"group_chat_created"`      // optional
	// SuperGroupChatCreated bool               `json:"supergroup_chat_created"` // optional
	// ChannelChatCreated    bool               `json:"channel_chat_created"`    // optional
	// MigrateToChatID       int64              `json:"migrate_to_chat_id"`      // optional
	// MigrateFromChatID     int64              `json:"migrate_from_chat_id"`    // optional
	// PinnedMessage         *Message           `json:"pinned_message"`          // optional
	// Invoice               *Invoice           `json:"invoice"`                 // optional
	// SuccessfulPayment     *SuccessfulPayment `json:"successful_payment"`      // optional
}

// User is a user on Telegram.
type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`     // optional
	UserName     string `json:"username"`      // optional
	LanguageCode string `json:"language_code"` // optional
	IsBot        bool   `json:"is_bot"`        // optional
}

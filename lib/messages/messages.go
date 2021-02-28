package messages

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// COOKIE NAME
const sessionName = "fmessages"
const sessionKey = "example-session-key" // TODO: THIS SHOULD NOT BE HARDCODED

// MESSAGE TYPES
type messageType string

const (
	Success		messageType = "success"
	Error		messageType = "error"
)

// GET A POINTER TO A sessions.CookieStore
func getCookieStore() *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(sessionKey))
}

// ADD A NEW MESSAGE INTO THE COOKIES STORAGE
func Set (c echo.Context, t messageType, value string) {

	name := string(t)
	session, _ := getCookieStore().Get(c.Request(), sessionName)
	session.AddFlash(value, name)
	session.Save(c.Request(), c.Response())
}

// GET A FLASHED MESSAGE FROM THE COOKIES STORAGE
func Get (c echo.Context, t messageType) []string {

	name := string(t)

	session, _ := getCookieStore().Get(c.Request(), sessionName)
	fm := session.Flashes(name)

	// IF WE HAVE MESSAGES
	if len(fm) > 0 {
		session.Save(c.Request(), c.Response())

		var flashes []string
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
		return flashes
	}

	return nil
}

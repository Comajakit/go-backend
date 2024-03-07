// config/session.go

package config

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var sessionStore sessions.Store

func init() {
	// Initialize the session store
	sessionStore = cookie.NewStore([]byte("secret")) // Replace "secret" with your own secret key
	// Optionally, set other session store configurations, such as MaxAge
	// sessionStore.Options(sessions.Options{MaxAge: ...})
}

// GetSessionStore returns the initialized session store
func GetSessionStore() sessions.Store {
	return sessionStore
}

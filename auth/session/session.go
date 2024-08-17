package auth

import (
	"sync"
	"time"
)

// Session structure
type Session struct {
	Login        string
	SessionToken string
	ExpiryDate   time.Time
}

// SessionManager structure
type SessionManager struct {
	sessions map[string]Session
	mutex    sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]Session),
	}
}

func (sm *SessionManager) ClearExpiredSessions() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for _, session := range sm.sessions {
		if time.Now().After(session.ExpiryDate) {
			delete(sm.sessions, session.SessionToken)
		}
	}
}

func (sm *SessionManager) ClearSession(sessionToken string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.sessions, sessionToken)
}

// CheckSession checks if a given session token is valid and not expired, and returns the associated user ID
func (sm *SessionManager) CheckSession(token string) (Session, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	session := sm.sessions[token]

	if time.Now().After(session.ExpiryDate) {
		return Session{}, false
	}

	return session, true
}

// CreateSession creates a new session for a user
func (sm *SessionManager) CreateSession(login string, sessionToken string) Session {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session := Session{
		Login:        login,
		SessionToken: sessionToken,
		ExpiryDate:   time.Now().Add(30 * time.Minute), // Session valid for 60 minutes
	}

	sm.sessions[sessionToken] = session
	return session
}

package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type SessionStore struct {
	mu sync.RWMutex
	sessions map[string]int64
}

func NewSessionStore() *SessionStore{
	return &SessionStore{
		sessions: make(map[string]int64),
	}
}

func (s *SessionStore) Create(userId int64) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	buffer := make([]byte,32)
	if _,err := rand.Read(buffer); err != nil {
		return "",err
	}
	token := hex.EncodeToString(buffer)
	
	s.sessions[token] = userId
	return  token,nil
}


func (s *SessionStore) Get(token string) (int64,bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userId,ok := s.sessions[token]
	return userId,ok
}
func (s *SessionStore) Delete(token string)  {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)

}



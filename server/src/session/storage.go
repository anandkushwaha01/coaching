package session

import (
	"time"
)

type SessionStore struct {
	Sid          string                 // unique session id
	TimeAccessed time.Time              // last access time
	Value        map[string]interface{} // session value stored inside
}

func (st *SessionStore) Set(key string, value interface{}) error {
	st.Value[key] = value
	return nil
}

func (st *SessionStore) Get(key string) interface{} {
	if v, ok := st.Value[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *SessionStore) Delete(key string) error {
	delete(st.Value, key)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.Sid
}

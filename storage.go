package main

import (
	"crypto/sha512"
	"fmt"
	"sync"
	"time"
)

type UserStorage interface {
	Set(user *User) error
	Del(key string) error
	Get(key string) *User
	List() []*User
}

type User struct {
	RemoteURL   string
	DeployKeys  []byte
	User        string
	Email       string
	ErrCount    int
	LastRunTime time.Time
}

func (u *User) Sign() string {
	s := fmt.Sprintf("%s-%s-%s", u.RemoteURL, u.User, u.Email)
	sum512 := sha512.Sum512([]byte(s))
	return fmt.Sprintf("%x", sum512)
}

// ---- 存储区

type MapStorage struct {
	sync.RWMutex
	Map map[string]*User
}

func NewMapStorage() *MapStorage {
	st := new(MapStorage)
	st.Map = make(map[string]*User)
	return st
}

func (m *MapStorage) Set(user *User) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Map[user.Sign()]; ok {
		return fmt.Errorf("已存在")
	}
	m.Map[user.Sign()] = user
	return nil
}

func (m *MapStorage) Del(key string) error {
	m.Lock()
	defer m.Unlock()
	delete(m.Map, key)
	return nil
}

func (m *MapStorage) Get(key string) *User {
	m.RLock()
	defer m.RUnlock()
	return m.Map[key]
}

func (m *MapStorage) List() []*User {
	m.RLock()
	defer m.RUnlock()
	list := make([]*User, 0)
	for _, user := range m.Map {
		list = append(list, user)
	}
	return list
}

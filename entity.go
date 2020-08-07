package main

import (
	"crypto/sha512"
	"fmt"
	"time"
)

type User struct {
	RemoteURL   string    `json:"remote_url"`
	DeployKeys  []byte    `json:"deploy_keys"`
	User        string    `json:"user"`
	Email       string    `json:"email"`
	ErrCount    int       `json:"err_count"`
	LastRunTime time.Time `json:"last_run_time"`
}

func (u *User) Sign() string {
	s := fmt.Sprintf("%s-%s-%s-%s", u.RemoteURL, u.User, u.Email, u.DeployKeys)
	sum512 := sha512.Sum512([]byte(s))
	return fmt.Sprintf("%x", sum512)
}

func (u *User) TOVO() *userResp {
	return &userResp{
		u.RemoteURL,
		u.User,
		u.Email,
		u.ErrCount,
		u.LastRunTime,
	}
}

type userParam struct {
	RemoteURL  string `form:"remote_url" json:"remote_url" binding:"required"`
	DeployKeys string `form:"deploy_keys" json:"deploy_keys" binding:"required"`
	User       string `form:"user" json:"user" binding:"required"`
	Email      string `form:"email" json:"email" binding:"required"`
}

func (p *userParam) TOUser() *User {
	return &User{
		RemoteURL:  p.RemoteURL,
		DeployKeys: []byte(p.DeployKeys),
		User:       p.User,
		Email:      p.Email,
	}
}

type userResp struct {
	RemoteURL   string    `json:"remote_url"`
	User        string    `json:"user"`
	Email       string    `json:"email"`
	ErrCount    int       `json:"err_count"`
	LastRunTime time.Time `json:"last_run_time"`
}

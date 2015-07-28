package main

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"tuluu.com/liut/staffio/backends/ldap"
	. "tuluu.com/liut/staffio/settings"
)

type User struct {
	Uid  string
	Name string
}

type Context struct {
	Session   *sessions.Session
	ResUrl    string
	User      *User
	NavSimple bool
	Referer   string
}

func (c *Context) Close() {
	ldap.CloseAll()
}

func NewContext(req *http.Request) (*Context, error) {
	sess, err := store.Get(req, Settings.Session.Name)
	sess.Options.Domain = Settings.Session.Domain
	sess.Options.HttpOnly = true
	var user *User
	if v, ok := sess.Values["user"]; ok {
		user = v.(*User)
	}
	referer := req.FormValue("referer")
	if referer == "" {
		referer = req.Referer()
	}
	ctx := &Context{
		Session: sess,
		ResUrl:  Settings.ResUrl,
		Referer: referer,
		User:    user,
	}
	if err != nil {
		log.Printf("new context error: %s", err)
		return ctx, err
	}

	return ctx, err
}

func init() {
	gob.Register(&User{})
}
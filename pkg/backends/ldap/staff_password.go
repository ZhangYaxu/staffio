package ldap

import (
	"github.com/go-ldap/ldap"
	"log"
)

func (s *LDAPStore) PasswordChange(uid, oldPasswd, newPasswd string) (err error) {
	for _, ls := range s.sources {
		err = ls.PasswordChange(uid, oldPasswd, newPasswd)
		if err != nil {
			log.Printf("PasswordChange at %s ERR: %s", ls.Addr, err)
		}
	}
	return
}

func (ls *ldapSource) PasswordChange(uid, oldPasswd, newPasswd string) error {
	userdn := ls.UDN(uid)
	err := ls.Bind(userdn, oldPasswd, true)
	if err != nil {
		return err
	}
	passwordModifyRequest := ldap.NewPasswordModifyRequest(userdn, oldPasswd, newPasswd)
	passwordModifyResponse, err := ls.c.PasswordModify(passwordModifyRequest)

	if err != nil {
		log.Printf("PasswordModify ERR: %s", err)
		return err
	}

	log.Printf("passwordModifyResponse: %v", passwordModifyResponse)
	return nil
}

func (s *LDAPStore) PasswordReset(uid, passwd string) (err error) {
	for _, ls := range s.sources {
		err = ls.PasswordReset(uid, passwd)
		if err != nil {
			log.Printf("PasswordReset at %s ERR: %s", ls.Addr, err)
		}
	}
	return
}

// password reset by administrator
func (ls *ldapSource) PasswordReset(uid, newPasswd string) error {
	err := ls.Bind(ls.BindDN, ls.Passwd, false)
	if err != nil {
		return err
	}
	dn := ls.UDN(uid)

	passwordModifyRequest := ldap.NewPasswordModifyRequest(dn, "", newPasswd)
	passwordModifyResponse, err := ls.c.PasswordModify(passwordModifyRequest)

	if err != nil {
		log.Printf("PasswordModify ERR: %s", err)
		return err
	}

	log.Printf("passwordModifyResponse: %v", passwordModifyResponse)
	return nil
}

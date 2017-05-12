package backends

import (
	"fmt"
	"log"
	"strings"

	"github.com/wealthworks/go-tencent-api/exmail"

	"lcgc/platform/staffio/models"
	. "lcgc/platform/staffio/settings"
)

func GetStaffFromExmail(email string) (*models.Staff, error) {
	user, err := exmail.GetUser(email)
	if err != nil {
		return nil, err
	}

	debug("got from exmail: %s", user)
	sn, gn := models.SplitName(user.Name)

	// log.Printf("got %q %q %q", user.Name, sn, gn)

	return &models.Staff{
		Uid:            strings.Split(user.Alias, "@")[0],
		Email:          user.Alias,
		CommonName:     user.Name,
		Surname:        sn,
		GivenName:      gn,
		EmployeeNumber: user.ExtId,
		EmployeeType:   user.Title,
		Mobile:         user.Mobile,
		Gender:         models.Gender(user.Gender),
	}, nil
}

func GetEmailAddress(uid string) string {
	return fmt.Sprintf("%s@%s", uid, Settings.EmailDomain)
}

func CheckMailUnseen(uid string) int {
	email := GetEmailAddress(uid)
	count, err := exmail.CountNewMail(email)
	if err != nil {
		log.Printf("check mail %s unseen ERR %s", uid, err)
	}
	return count
}

func GetMailEntryUrl(uid string) string {
	email := GetEmailAddress(uid)
	str, err := exmail.GetLoginURL(email)
	if err != nil {
		log.Printf("get login url of %s, ERR %s", uid, err)
	}
	return str
}

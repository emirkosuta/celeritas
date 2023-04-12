package session

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	CookieSecure   string
	DBPool         *sql.DB
}

func (c *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	// how long should sessions last?
	minutes, err := strconv.Atoi(c.CookieLifetime)
	if err != nil {
		minutes = 60
	}

	// should cookies perist?
	if strings.ToLower(c.CookiePersist) == "true" {
		persist = true
	}

	// must cookies be secure?
	if strings.ToLower(c.CookieSecure) == "true" {
		secure = true
	}

	// create session
	sess := scs.New()
	sess.Lifetime = time.Duration(minutes) * time.Minute
	sess.Cookie.Persist = persist
	sess.Cookie.Name = c.CookieName
	sess.Cookie.Secure = secure
	sess.Cookie.Domain = c.CookieDomain
	sess.Cookie.SameSite = http.SameSiteLaxMode

	// which session store?
	switch strings.ToLower(c.SessionType) {
	case "redis":

	case "mysql", "mariadb":
		sess.Store = mysqlstore.New(c.DBPool)
	case "postgres", "postgresql":
		sess.Store = postgresstore.New(c.DBPool)
	default:
		// cookie
	}

	return sess
}

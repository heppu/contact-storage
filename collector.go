package collector

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Collector struct {
	db    *sql.DB
	dbmap *gorp.DbMap
}

func Run(listenAddr, dbUrl string) (err error) {
	// Init DB
	connStr, err := pq.ParseURL(dbUrl)
	if err != nil {
		return
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return
	}

	c := &Collector{
		db: db,
		dbmap: &gorp.DbMap{
			Db:      db,
			Dialect: gorp.PostgresDialect{},
		},
	}

	c.dbmap.AddTableWithName(Contact{}, "contacts").SetKeys(true, "ID").AddIndex("email_index", "B-TREE", []string{"email"}).SetUnique(true)
	if err = c.dbmap.CreateTablesIfNotExists(); err != nil {
		return
	}

	// Start server
	router := httprouter.New()
	router.POST("/contacts", c.addContact)
	return http.ListenAndServe(listenAddr, router)
}

func (c *Collector) addContact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	contact := Contact{}
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		log.Errorln("Could not decode request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contact.CreatedAt = time.Now()
	if err := c.dbmap.Insert(&contact); err != nil {
		log.Errorln("Could not insert contact into db:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type Contact struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

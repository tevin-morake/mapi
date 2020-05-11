package tools

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/CaboodleData/gotools/file"
	_ "github.com/lib/pq"
)

var (
	dbhost   = os.Getenv(DB_HOST)
	dbdriver = os.Getenv(DB_DRIVER)
	dbport   = os.Getenv(DB_PORT)
	dbuser   = os.Getenv(DB_USER)
	dbpass   = os.Getenv(DB_PASSWORD)
	dbname = os.Getenv(DB_NAME)
	stdout   = false
)

type Tools struct {
	Info       *log.Logger
	Error      *log.Logger
	PostgresDb *sql.DB
}

//NewTools initializes an instance of tools we can use throughout the app
func NewTools(stdout bool, filename string, logname string) (*Tools, error) {
	var t Tools
	var err error	
	t.Info, t.Error = InitLogs(stdout, "logs", "mapi")//Initialize log file on the server
	return &t, err
}

//InitLogs initializes loggers for the application logs
func InitLogs(stdout bool, logFolder string, prefix string) (*log.Logger, *log.Logger) {
	var handler io.Writer
	var err error

	if stdout { // STDOUT logs to console, else we log to a file
		handler = os.Stdout
	} else {
		// Create folder and a log file name indicating todays date
		os.Mkdir(logFolder, 0777)
		logName := filepath.Join(logFolder, prefix+"_"+time.Now().Format("20060102")+".log")
		handler, err = os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("Error opening logs " + err.Error())
		}
		//Keep logfolder clean by deleting logs older than 30 days
		_ = CleanFolder(logFolder, 10)
	}

	infoLog := log.New(handler, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(handler, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return infoLog, errorLog
	var handler io.Writer
	var err error
	// STDOUT logs to console, else we log to a file
	if stdout {
		handler = os.Stdout
	} else {
		// Create folder and a log file name indicating todays date
		os.Mkdir(logFolder, 0777)
		logName := filepath.Join(logFolder, prefix+"_"+time.Now().Format("20060102")+".log")
		handler, err = os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("Error opening logs " + err.Error())
		}
		//Keep logfolder clean by deleting logs older than 30 days
		_ = CleanFolder(logFolder, 10)
	}

	infoLog := log.New(handler, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(handler, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return infoLog, errorLog
}

//CleanFolder cleans the log folders after a given interval
func CleanFolder(root string, daystokeep int) error {
	dateToKeep := time.Now().AddDate(0, 0, daystokeep*-1)
	fn := func(path string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			return nil
		}
		if f.ModTime().Before(dateToKeep) {
			err = os.Remove(filepath.Join(root, f.Name()))
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := filepath.Walk(root, fn)
	if err != nil {
		return err
	}
	return nil
	dateToKeep := time.Now().AddDate(0, 0, daystokeep*-1)
	fn := func(path string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			return nil
		}
		if f.ModTime().Before(dateToKeep) {
			err = os.Remove(filepath.Join(root, f.Name()))
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := filepath.Walk(root, fn)
	if err != nil {
		return err
	}
	return nil
}

//OpenDB opens connection to pgsql db
func (t *Tools) OpenDB(w http.ResponseWriter) error {
	if t.PostgresDb == nil {
		pgsqlConnString := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",dbhost,dbport, dbuser, dbpass, dbname)
		t.PostgresDb, err := sql.Open(dbdriver)
		if err != nil {
			return err
		}
		defer t.PostgresDb.Close()

		if err := t.PostgresDb.Ping(); err != nil {
			return err
		}
		fmt.Println("Connection to db opened !")
	}

}

//ExecuteSQL runs a prepared sql statement
func (t *Tools) ExecuteSQL(sqlStatement string, values interface{}) error {
	// Get pgsql transactions
	pgsqlTx, err := t.PostgresDb.Begin()
	if err != nil {
		log.Printf("Error getting PostgreSQL transaction|%s", err.Error())
		return err
	}

	//Prepare sql statement for db insertion
	pgsqlStatement, err := t.PostgresDb.Prepare(sqlStatement)
	if err != nil {
		pgsqlTx.Rollback()
		log.Printf("Error preparing postgres %s", err.Error())
		return err
	}

	//execute statement. rollback if err occurs
	if _, err = pgsqlTx.Exec(values...); err != nil {
		log.Printf("Error executing prepared pgsql statement : %s", err.Error())
		pgsqlTx.Rollback()
		return err
	}

	// Commit transaction to pgsql db
	pgsqlTx.Commit()
	return nil

}

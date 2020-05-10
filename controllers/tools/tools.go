package tools

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

//InitLogs initializes loggers for the application logs
func InitLogs(stdout bool, logFolder string, prefix string) (*log.Logger, *log.Logger) {
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
}

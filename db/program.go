package db

import (
	"encoding/json"

	"github.com/VolticFroogo/Launcher-Server/model"
)

// GetProgram will get a program from the database given an ID.
func GetProgram(id string) (program model.Program, err error) {
	rows, err := db.Query("SELECT name, versions, path, exceptions FROM programs WHERE id=?", id)
	if err != nil {
		return
	}

	defer rows.Close()

	program.ID = id

	if rows.Next() {
		var versionsJSON, exceptionsJSON string

		err = rows.Scan(&program.Name, &versionsJSON, &program.Path, &exceptionsJSON)
		if err != nil {
			return
		}

		err = json.Unmarshal([]byte(versionsJSON), &program.Versions)
		if err != nil {
			return
		}

		err = json.Unmarshal([]byte(exceptionsJSON), &program.Exceptions)
	}

	return
}

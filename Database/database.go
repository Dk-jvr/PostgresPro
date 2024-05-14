package Database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"postgres/Models"
	"sync"
)

var (
	db    *sql.DB
	mutex sync.Mutex
)

func ConnectDatabase() *sql.DB {
	var err error

	const connectionString = `user=dbuser dbname=pg_commands password=password host=db port=5432 sslmode=disable`
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func AddScript(script Models.ScriptData) (int64, error) {
	var scriptId int64
	const query = `INSERT INTO Scripts(script, script_type) VALUES ($1, $2) RETURNING script_id`
	mutex.Lock()
	defer mutex.Unlock()
	err := db.QueryRow(query, script.Script, script.Type).Scan(&scriptId)
	return scriptId, err
}

func GetScript(id int64) (Models.ScriptData, error) {
	var data Models.ScriptData
	const query = `SELECT script_id, script, script_type FROM Scripts WHERE script_id = $1`
	mutex.Lock()
	defer mutex.Unlock()
	err := db.QueryRow(query, id).Scan(&data.Id, &data.Script, &data.Type)
	return data, err
}

func GetScripts() ([]Models.ScriptData, error) {
	var (
		rows    *sql.Rows
		scripts []Models.ScriptData
		err     error
	)
	const query = `SELECT script_id, script, script_type FROM Scripts`
	mutex.Lock()
	defer mutex.Unlock()
	rows, err = db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		data := new(Models.ScriptData)
		err = rows.Scan(&data.Id, &data.Script, &data.Type)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, *data)
	}
	return scripts, nil
}

func InsertScriptData(data Models.ScriptData, output string) error {
	const queryInsertData = `INSERT INTO OutPutScripts(script_id, script, script_type, output, output_time) VALUES ($1, $2, $3, $4, NOW())`
	mutex.Lock()
	defer mutex.Unlock()
	_, err := db.Exec(queryInsertData, data.Id, data.Script, data.Type, output)
	return err
}

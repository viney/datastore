package datastore

import (
	"fmt"
	"testing"
	"time"
)

var userTableSql string = `
BEGIN
    BEGIN
         EXECUTE IMMEDIATE 'DROP TABLE user_profile';
    EXCEPTION
         WHEN OTHERS THEN
                IF SQLCODE != -942 THEN
                     RAISE;
                END IF;
    END;
    EXECUTE IMMEDIATE 'CREATE TABLE user_profile (userid NUMBER(10) PRIMARY KEY, name VARCHAR(20) NOT NULL, created VARCHAR(20) NOT NULL)';
END;
`

func TestDB(t *testing.T) {
	// open
	db, err := Open()
	checkErr(t, err)

	defer Close(db)

	_, err = db.Exec(userTableSql)
	checkErr(t, err)

	// insert
	tx, err := db.Begin()
	checkErr(t, err)

	insertSql := `insert into user_profile(userid,name,created) values(:1,:2,:3)`
	stmt, err := tx.Prepare(insertSql)
	checkErr(t, err)

	_, err = stmt.Exec(1, "viney", time.Now().Format("2006-01-02 15:04:05"))
	checkErr(t, err)

	// update
	updateSql := `update user_profile set name=:1 where userid=1`
	stmt, err = tx.Prepare(updateSql)
	checkErr(t, err)
	_, err = stmt.Exec("中国人")
	checkErr(t, err)

	// select
	querySql := `select userid,name,created from user_profile`
	rows, err := tx.Query(querySql)

	type user struct {
		userid  int
		name    string
		created string
	}

	var u = &user{}
	for rows.Next() {
		err = rows.Scan(
			&u.userid,
			&u.name,
			&u.created)
		checkErr(t, err)
	}
	rows.Close()

	fmt.Println(*u)

	// delete
	/*deleteSql := `delete from user_profile where userid=1`
	_, err = db.Exec(deleteSql)
	checkErr(err)*/

	err = tx.Commit()
	checkErr(t, err)
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Error("oracle test: ", err.Error())
		return
	}
}

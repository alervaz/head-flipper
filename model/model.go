package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	id string
}

type User struct {
	username string
	points   int
	serverId string
}

var DB, err = sql.Open("sqlite3", "model/db.sqlite")

func PingDB() error {
	return DB.Ping()
}

func CreateServer(serverId string) error {
	_, err := DB.Exec("INSERT INTO servers (id) VALUES ($1)", serverId)
	return err
}

func CreateUser(username string, serverId string) error {
	_, err := DB.Exec(
		"INSERT INTO users (username, points, serverId) VALUES ($1, $2, $3)",
		username,
		0,
		serverId,
	)
	return err
}

func UpdatePoints(username string, points int) error {
	_, err := DB.Exec("UPDATE users SET points = ? WHERE username = ?", points, username)
	return err
}

func Gamble(username string, serverId string, points int) error {
	row := DB.QueryRow("SELECT * FROM servers WHERE id = ?", serverId)
	var server Server
	if row.Scan(&server.id) != nil {
		err := CreateServer(serverId)
		if err != nil {
			return err
		}
	}

	row = DB.QueryRow("SELECT * FROM users WHERE username = ?", username)
	var user User
	if row.Scan(&user.username, &user.points, &user.serverId) != nil {
		err := CreateUser(username, serverId)
		if err != nil {
			return err
		}
	}

	err := UpdatePoints(username, points)
	return err
}

func GetPoints(username string) (int, error) {
	row := DB.QueryRow("SELECT points FROM users WHERE username = ?", username)
	var user User

	err := row.Scan(&user.points)
	return user.points, err
}

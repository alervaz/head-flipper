package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	id string
}

type User struct {
	id       uint
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

func UpdatePoints(username string, serverId string, points int) error {
	_, err := DB.Exec("UPDATE users SET points = ? WHERE username = ? AND serverId = ?", points, username, serverId)
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

	row = DB.QueryRow("SELECT * FROM users WHERE username = ? AND serverId = ?", username, serverId)
	var user User
	if row.Scan(&user.id, &user.username, &user.points, &user.serverId) != nil {
		err := CreateUser(username, serverId)
		if err != nil {
			return err
		}
	}

	err := UpdatePoints(username, serverId, points)
	return err
}

func GetPoints(username string, serverId string) (int, error) {
	row := DB.QueryRow("SELECT * FROM servers WHERE id = ?", serverId)
	var server Server
	if err := row.Scan(&server.id); err != nil {
		err := CreateServer(serverId)
		if err != nil {
			return 0, err
		}
	}

	row = DB.QueryRow("SELECT * FROM users WHERE username = ? AND serverId = ?", username, serverId)
	var user User
	if err := row.Scan(&user.id, &user.username, &user.points, &user.serverId); err != nil {
		err := CreateUser(username, serverId)
		if err != nil {
			return 0, err
		}
	}

	return user.points, err
}

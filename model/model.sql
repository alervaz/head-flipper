DROP TABLE IF EXISTS servers;
DROP TABLE IF EXISTS users;

CREATE TABLE servers (
  id TEXT PRIMARY KEY UNIQUE 
);

CREATE TABLE users (
  username TEXT PRIMARY KEY UNIQUE,
  points INTEGER,
  serverId TEXT,
  FOREIGN KEY(serverId) REFERENCES servers(id)
)



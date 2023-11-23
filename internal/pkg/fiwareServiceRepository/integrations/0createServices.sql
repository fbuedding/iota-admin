CREATE TABLE IF NOT EXISTS"services" (
	"name"	TEXT NOT NULL,
	"id"	TEXT NOT NULL UNIQUE,
	"created_at"	DATETIME NOT NULL,
	"updated_at"	DATETIME NOT NULL,
	PRIMARY KEY("id")
)
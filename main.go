package main

import (
	"database/sql"
	"log"
	"net/http"

	"sql-injection-subtle/internal/database"
	"sql-injection-subtle/internal/handler"
	"sql-injection-subtle/internal/repository"
	"sql-injection-subtle/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:data.db?mode=memory&cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.InitSchema(db); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userH := handler.NewUserHandler(userSvc)

	http.HandleFunc("/users/search", userH.Search)
	http.HandleFunc("/users", userH.List)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"SQL injection subtle API","endpoints":["GET /users/search?q=","GET /users?username=","GET /users?sort=&order="]}`))
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/martini-contrib/binding"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/render"
)

type User struct {
	Mobile   string `json:"mobile"`
	Username string `json:"username"`
}

func main() {
	db, err := sql.Open("mysql", "root:ChangeMe@/monkey")
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	defer db.Close()

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Logger())

	m.Get("/products/:hot", func(params martini.Params, r render.Render, request *http.Request) {
		query := request.URL.Query()
		var limit, offset string
		if query.Get("limit") != "" {
			limit = " LIMIT " + query.Get("limit")
		}

		if query.Get("offset") != "" {
			offset = " OFFSET " + query.Get("offset")
		}

		sqlCommand := "SELECT id, name FROM monkey_company" + limit + offset
		fmt.Println(sqlCommand)
		res, _ := db.Query("select id, name from monkey_company")
		if err != nil {
			fmt.Println(err)
		}

		for res.Next() {
			var name string
			var id int
			if err := res.Scan(&id, &name); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s is %d\n", name, id)
		}
		fmt.Println(params["hot"])
		r.JSON(200, map[string]interface{}{"hello": "world"})
	})

	m.Post("/register", binding.Json(User{}), func(user User, render render.Render, request *http.Request) {

		render.JSON(200, &user)
	})
	m.Run()
}

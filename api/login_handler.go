package api

import (
	"fmt"
	"lifeofsems-go/types"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) HandleLogin(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	c := GetSessionCookie(req)
	http.SetCookie(w, c)
	user := s.GetUser(w, req)

	data := struct {
		Header types.Header
		Text   string
	}{
		Header: types.Header{
			Navigation: s.BuildNavigationItems(req),
			User:       "",
		},
		Text: "",
	}

	if req.Method == http.MethodGet {
		if user == nil {
			s.renderTemplate(w, req, "login", data)
			return
		}

		http.Redirect(w, req, "/admin", http.StatusSeeOther)
	} else if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}

		data.Text = "Invalid username or password."

		username := req.Form.Get("username")
		password := req.Form.Get("password")

		user, err = s.store.GetUserBy(map[string]string{"username": username})
		if err != nil {
			data.Text = "User does not exist."
			s.renderTemplate(w, req, "login", data)
			return
		}

		fmt.Println(password, string(user.Password))

		pw, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		fmt.Println(string(pw))

		if bcrypt.CompareHashAndPassword(user.Password, []byte(password)) != nil {
			s.renderTemplate(w, req, "login", data)
			return
		}

		s.store.CreateSession(c.Value, username)
		http.Redirect(w, req, "/admin", http.StatusSeeOther)
	}
}

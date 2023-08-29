package api

import (
	"fmt"
	"lifeofsems-go/types"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) HandleBlogPage(w http.ResponseWriter, req *http.Request) {
	tokens := strings.Split(req.URL.Path, "/")

	if len(tokens) < 3 {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	// POST on blog/create
	if tokens[2] == "create" {
		if req.Method == http.MethodPost {
			s.CreatePost(w, req)
			return
		} else {
			s.HandleErrorPage(w, req, http.StatusMethodNotAllowed)
			return
		}
	}

	// GET, PUT, DELETE on blog/{postId}
	postId, err := strconv.Atoi(tokens[2])
	if err != nil {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	if req.Method == http.MethodGet {
		s.ViewPost(w, req, postId)
	} else if req.Method == http.MethodPut {
		fmt.Println("Method put on blog/{:d}")
	} else if req.Method == http.MethodDelete {
		fmt.Println("Method delete on blog/{:d}")
	} else {
		s.HandleErrorPage(w, req, http.StatusMethodNotAllowed)
	}
}

func (s *Server) CreatePost(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (s *Server) ViewPost(w http.ResponseWriter, req *http.Request, postId int) {
	blogPost, err := s.store.GetPost(postId)
	if err != nil {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	data := struct {
		Header types.Header
		Post   *types.BlogPost
	}{
		Header: types.Header{
			Navigation: s.BuildNavigationItems(req),
			User:       "",
		},
		Post: blogPost,
	}

	w.Header().Add("Content-Type", "text/html")
	s.renderTemplate(w, req, "blog-post", data)
}

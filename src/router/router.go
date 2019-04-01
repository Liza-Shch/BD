package router

import (
	"../handlers"
	"github.com/gorilla/mux"
)

//TODO: переписать подключение с ипользованием методов get post + разделение самих хендлеров по этоиу принципу

func Routing(r *mux.Router) {
	r.HandleFunc("/api/user/{nickname}/create", handlers.CreateUser)
	r.HandleFunc("/api/user/{nickname}/profile", handlers.ProfileUser)
	r.HandleFunc("/api/forum/create", handlers.CreateForum)
	r.HandleFunc("/api/forum/{slug}/details", handlers.GetForum)
	r.HandleFunc("/api/forum/{slug}/create", handlers.CreateThread)
	r.HandleFunc("/api/thread/{slug_or_id}/details", handlers.ThreadDetails)
	r.HandleFunc("/api/thread/{slug_or_id}/vote", handlers.VoteThread)
	r.HandleFunc("/api/forum/{slug}/threads", handlers.GetThreads)
	r.HandleFunc("/api/thread/{slug_or_id}/create", handlers.CreatePost)
	r.HandleFunc("/api/forum/{slug}/users", handlers.GetUsersByForum)
	r.HandleFunc("/api/post/{id}/details", handlers.PostDetails)
	r.HandleFunc("/api/service/status", handlers.ServiceStatus)
	r.HandleFunc("/api/service/clear", handlers.ServiceClear)
	r.HandleFunc("/api/thread/{slug_or_id}/posts", handlers.GetPosts)
}

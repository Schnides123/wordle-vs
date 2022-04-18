package endpoints

import "github.com/gorilla/mux"

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/wordle", WordleHandler)
	r.HandleFunc("/blitz/", MatchmakingHandler)
	r.HandleFunc("/blitz/challenge", ChallengeLinkHandler)
	r.HandleFunc("/blitz/{gameID}", BlitzHandler)
}

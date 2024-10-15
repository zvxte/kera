package server

type Server struct {
	authRouter     Router
	usersRouter    Router
	sessionsRouter Router
	habitsRouter   Router
}

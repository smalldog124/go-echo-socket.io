package main

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Middleware for handling CORS
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE,OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			return next(c)
		}
	})

	server := socketio.NewServer(nil)

	//Register the WebSocket endpoint
	e.GET("/socket.io/", func(c echo.Context) error {
		server.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Handle WebSocket connections
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("Client socket connected:", s.ID())
		return nil
	})

	// Handle "order" events from clients
	server.OnEvent("/", "order", func(s socketio.Conn, msg string) {
		log.Println("Order received:", msg)
		s.Emit("90", "lek->com")

		// Emitting a "shop" event to all connected clients
		//server.BroadcastToNamespace("/", "shop", msg)
	})

	// Handle WebSocket disconnections
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("Client socket disconnected:", s.ID(), reason)
	})

	go server.Serve()
	defer server.Close()

	e.HideBanner = true

	e.Static("/", "../web")
	e.Any("/socket.io/", func(c echo.Context) error {
		server.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Start the Echo server
	e.Logger.Fatal(e.Start(":3000"))
}

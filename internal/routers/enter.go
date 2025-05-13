package routers

type RouterGroup struct {
	Showtime ShowtimeRouter
}

var ShowtimeServiceRouterGroup = new(RouterGroup)

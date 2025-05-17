package routers

type RouterGroup struct {
	Showtime     ShowtimeRouter
	ShowtimeSeat ShowtimeSeatRouter
}

var ShowtimeServiceRouterGroup = new(RouterGroup)

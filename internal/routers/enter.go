package routers

type RouterGroup struct {
	Showtime     ShowtimeRouter
	ShowtimeSeat ShowtimeSeatRouter
	Cinema       CinemaRouter
}

var ShowtimeServiceRouterGroup = new(RouterGroup)

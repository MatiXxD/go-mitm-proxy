package webapi

import "github.com/MatiXxD/go-mitm-proxy/internal/delivery/request"

func (s *Server) BindRoutes(rd *request.RequestDelivery) {
	s.echo.GET("/requests", rd.GetRequestsInfo())
	s.echo.GET("/requests/:id", rd.GetRequestById())
	s.echo.GET("/repeat/:id", rd.RepeatRequest())
	s.echo.GET("/scan/:id", rd.ScanRequest())
}

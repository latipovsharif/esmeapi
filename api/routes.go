package api

// GetRoutes return routes for api
func (s *Server) getRoutes() {
	api := s.e.Group("/api/v1")
	api.GET("/balance/", s.balanceView)
	api.POST("/sms/send/", s.sendSMSView)
	api.GET("/sms/:id/", s.smsDetailView)
}

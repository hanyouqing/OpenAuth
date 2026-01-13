package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/hanyouqing/openauth/internal/services"
	"github.com/sirupsen/logrus"
)

type CASHandler struct {
	service *services.CASService
	logger  *logrus.Logger
}

func NewCASHandler(service *services.CASService, logger *logrus.Logger) *CASHandler {
	return &CASHandler{service: service, logger: logger}
}

// CASLogin handles CAS protocol login
// @Summary CAS Login
// @Description CAS (Central Authentication Service) protocol login endpoint
// @Tags sso
// @Produce html
// @Param service query string false "Service URL to redirect after login"
// @Success 200 "Login page or redirect to service"
// @Router /cas/login [get]
func (h *CASHandler) CASLogin(c *gin.Context) {
	h.service.CASLogin(c)
}

// CASValidate handles CAS 1.0 ticket validation
// @Summary CAS 1.0 Validate
// @Description CAS 1.0 protocol ticket validation endpoint
// @Tags sso
// @Produce text/plain
// @Param ticket query string true "Service ticket"
// @Param service query string true "Service URL"
// @Success 200 "yes\nusername or no\n"
// @Router /cas/validate [get]
func (h *CASHandler) CASValidate(c *gin.Context) {
	h.service.CASValidate(c)
}

// CASServiceValidate handles CAS 2.0 service ticket validation
// @Summary CAS 2.0 Service Validate
// @Description CAS 2.0 protocol service ticket validation endpoint (XML format)
// @Tags sso
// @Produce application/xml
// @Param ticket query string true "Service ticket"
// @Param service query string true "Service URL"
// @Success 200 "CAS 2.0 XML response"
// @Failure 400 "Invalid ticket"
// @Router /cas/serviceValidate [get]
func (h *CASHandler) CASServiceValidate(c *gin.Context) {
	h.service.CASServiceValidate(c)
}

// CASLogout handles CAS protocol logout
// @Summary CAS Logout
// @Description CAS protocol logout endpoint
// @Tags sso
// @Produce html
// @Param service query string false "Service URL to redirect after logout"
// @Success 200 "Logout page or redirect"
// @Router /cas/logout [get]
func (h *CASHandler) CASLogout(c *gin.Context) {
	h.service.CASLogout(c)
}

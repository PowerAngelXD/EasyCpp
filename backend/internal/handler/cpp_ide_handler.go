package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"easycpp/backend/internal/dto"
	"easycpp/backend/internal/service"
	"easycpp/backend/pkg/response"
)

// CPPIdeHandler 提供 C++ 在线编译执行接口。
type CPPIdeHandler struct {
	cppIdeService service.CPPIdeService
}

func NewCPPIdeHandler(cppIdeService service.CPPIdeService) *CPPIdeHandler {
	return &CPPIdeHandler{cppIdeService: cppIdeService}
}

func (h *CPPIdeHandler) RunCPP(c *gin.Context) {
	var req dto.RunCPPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.cppIdeService.RunCPP(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to run cpp code")
		return
	}

	response.Success(c, http.StatusOK, result)
}

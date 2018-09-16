package browse

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/portainer/agent"
	"github.com/portainer/agent/filesystem"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

type browseRenamePayload struct {
	CurrentFilePath string
	NewFilePath     string
}

func (payload *browseRenamePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.CurrentFilePath) {
		return agent.Error("Current file path is invalid")
	}
	if govalidator.IsNull(payload.NewFilePath) {
		return agent.Error("New file path is invalid")
	}
	return nil
}

func (handler *Handler) browseRename(rw http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	volumeID, err := request.RetrieveRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid volume identifier route variable", err}
	}

	var payload browseRenamePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	err = filesystem.RenameFileInsideVolume(volumeID, payload.CurrentFilePath, payload.NewFilePath)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to rename file", err}
	}

	return response.Empty(rw)
}

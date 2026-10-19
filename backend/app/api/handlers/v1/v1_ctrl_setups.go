package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/data/repo"
	"github.com/sysadminsmedia/homebox/backend/internal/web/adapters"
)

// HandleSetupsGetAll godoc
//
//	@Summary	Get All Setups
//	@Tags		Setups
//	@Produce	json
//	@Success	200	{object}	Results[repo.SetupOut]
//	@Router		/v1/setups [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupsGetAll() errchain.HandlerFunc {
	fn := func(r *http.Request, _ struct{}) ([]repo.SetupOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.GetAll(auth, auth.GID)
	}

	return adapters.Query(fn, http.StatusOK)
}

// HandleSetupsCreate godoc
//
//	@Summary	Create Setup
//	@Tags		Setups
//	@Produce	json
//	@Param		payload	body	repo.SetupCreate	true	"Setup Data"
//	@Success	201		{object}	repo.SetupOut
//	@Router		/v1/setups [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupsCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, body repo.SetupCreate) (repo.SetupOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.Create(auth, auth.GID, body)
	}

	return adapters.Action(fn, http.StatusCreated)
}

// HandleSetupGet godoc
//
//	@Summary	Get Setup
//	@Tags		Setups
//	@Produce	json
//	@Param		id	path	string	true	"Setup ID"
//	@Success	200	{object}	repo.SetupOut
//	@Router		/v1/setups/{id} [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) (repo.SetupOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.GetOne(auth, auth.GID, ID)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

// HandleSetupUpdate godoc
//
//	@Summary	Update Setup
//	@Tags		Setups
//	@Produce	json
//	@Param		id		path	string				true	"Setup ID"
//	@Param		payload	body	repo.SetupUpdate	true	"Setup Data"
//	@Success	200		{object}	repo.SetupOut
//	@Router		/v1/setups/{id} [PUT]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupUpdate() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, body repo.SetupUpdate) (repo.SetupOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.Update(auth, auth.GID, ID, body)
	}

	return adapters.ActionID("id", fn, http.StatusOK)
}

// HandleSetupDelete godoc
//
//	@Summary	Delete Setup
//	@Tags		Setups
//	@Param		id	path	string	true	"Setup ID"
//	@Success	204
//	@Router		/v1/setups/{id} [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupDelete() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) (any, error) {
		auth := services.NewContext(r.Context())
		return nil, ctrl.repo.Setups.Delete(auth, auth.GID, ID)
	}

	return adapters.ActionID("id", fn, http.StatusNoContent)
}

// HandleSetupPartsGet godoc
//
//	@Summary	Get Setup Parts (mounting history)
//	@Tags		Setups
//	@Produce	json
//	@Param		id	path	string	true	"Setup ID"
//	@Success	200	{object}	Results[repo.SetupPartOut]
//	@Router		/v1/setups/{id}/parts [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupPartsGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) ([]repo.SetupPartOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.GetParts(auth, auth.GID, ID, false)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

// HandleSetupPartsCurrentGet godoc
//
//	@Summary	Get Currently Mounted Parts
//	@Tags		Setups
//	@Produce	json
//	@Param		id	path	string	true	"Setup ID"
//	@Success	200	{object}	Results[repo.SetupPartOut]
//	@Router		/v1/setups/{id}/parts/current [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupPartsCurrentGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) ([]repo.SetupPartOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.GetParts(auth, auth.GID, ID, true)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

// HandleSetupPartAdd godoc
//
//	@Summary	Mount a Part on a Setup
//	@Tags		Setups
//	@Produce	json
//	@Param		id		path	string					true	"Setup ID"
//	@Param		payload	body	repo.SetupPartCreate	true	"Part Data"
//	@Success	201		{object}	repo.SetupPartOut
//	@Router		/v1/setups/{id}/parts [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupPartAdd() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, body repo.SetupPartCreate) (repo.SetupPartOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.AddPart(auth, auth.GID, ID, body)
	}

	return adapters.ActionID("id", fn, http.StatusCreated)
}

// HandleSetupPartRemove godoc
//
//	@Summary	Remove a Part from a Setup (stamps removedDate)
//	@Tags		Setups
//	@Param		id			path	string	true	"Setup ID"
//	@Param		part_id		path	string	true	"Setup Part ID"
//	@Success	204
//	@Router		/v1/setups/{id}/parts/{part_id} [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupPartRemove() errchain.HandlerFunc {
	fn := func(r *http.Request, partID uuid.UUID, _ struct{}) (any, error) {
		auth := services.NewContext(r.Context())
		return nil, ctrl.repo.Setups.RemovePart(auth, auth.GID, partID)
	}

	return adapters.ActionID("part_id", fn, http.StatusNoContent)
}

// HandleSetupPartDelete godoc
//
//	@Summary	Permanently delete a Setup Part assignment
//	@Tags		Setups
//	@Param		id			path	string	true	"Setup ID"
//	@Param		part_id		path	string	true	"Setup Part ID"
//	@Success	204
//	@Router		/v1/setups/{id}/parts/{part_id}/purge [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupPartDelete() errchain.HandlerFunc {
	fn := func(r *http.Request, partID uuid.UUID, _ struct{}) (any, error) {
		auth := services.NewContext(r.Context())
		return nil, ctrl.repo.Setups.DeletePart(auth, auth.GID, partID)
	}

	return adapters.ActionID("part_id", fn, http.StatusNoContent)
}

// HandleSetupSessionsGet godoc
//
//	@Summary	Get Usage Sessions for a Setup
//	@Tags		Setups
//	@Produce	json
//	@Param		id	path	string	true	"Setup ID"
//	@Success	200	{object}	Results[repo.UsageSessionOut]
//	@Router		/v1/setups/{id}/sessions [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupSessionsGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) ([]repo.UsageSessionOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.GetSessions(auth, auth.GID, ID)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

// HandleSetupSessionCreate godoc
//
//	@Summary	Log a Usage Session
//	@Tags		Setups
//	@Produce	json
//	@Param		id		path	string					true	"Setup ID"
//	@Param		payload	body	repo.UsageSessionCreate	true	"Session Data"
//	@Success	201		{object}	repo.UsageSessionOut
//	@Router		/v1/setups/{id}/sessions [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupSessionCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, body repo.UsageSessionCreate) (repo.UsageSessionOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.CreateSession(auth, auth.GID, ID, body)
	}

	return adapters.ActionID("id", fn, http.StatusCreated)
}

// HandleSetupSessionUpdate godoc
//
//	@Summary	Update a Usage Session
//	@Tags		Setups
//	@Produce	json
//	@Param		id			path	string					true	"Setup ID"
//	@Param		session_id	path	string					true	"Session ID"
//	@Param		payload		body	repo.UsageSessionUpdate	true	"Session Data"
//	@Success	200		{object}	repo.UsageSessionOut
//	@Router		/v1/setups/{id}/sessions/{session_id} [PUT]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupSessionUpdate() errchain.HandlerFunc {
	fn := func(r *http.Request, sessionID uuid.UUID, body repo.UsageSessionUpdate) (repo.UsageSessionOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.UpdateSession(auth, auth.GID, sessionID, body)
	}

	return adapters.ActionID("session_id", fn, http.StatusOK)
}

// HandleSetupSessionDelete godoc
//
//	@Summary	Delete a Usage Session
//	@Tags		Setups
//	@Param		id			path	string	true	"Setup ID"
//	@Param		session_id	path	string	true	"Session ID"
//	@Success	204
//	@Router		/v1/setups/{id}/sessions/{session_id} [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleSetupSessionDelete() errchain.HandlerFunc {
	fn := func(r *http.Request, sessionID uuid.UUID, _ struct{}) (any, error) {
		auth := services.NewContext(r.Context())
		return nil, ctrl.repo.Setups.DeleteSession(auth, auth.GID, sessionID)
	}

	return adapters.ActionID("session_id", fn, http.StatusNoContent)
}

// HandleEntitySetupHistoryGet godoc
//
//	@Summary	Get the mounting history of a part across all setups
//	@Tags		Setups
//	@Produce	json
//	@Param		id	path	string	true	"Entity (Part) ID"
//	@Success	200	{object}	Results[repo.SetupPartOut]
//	@Router		/v1/entities/{id}/setup-history [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleEntitySetupHistoryGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) ([]repo.SetupPartOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Setups.GetPartsByEntity(auth, auth.GID, ID)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

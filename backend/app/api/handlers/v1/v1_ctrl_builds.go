package v1

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/data/repo"
	"github.com/sysadminsmedia/homebox/backend/internal/sys/validate"
	"github.com/sysadminsmedia/homebox/backend/internal/web/adapters"
)

// HandleBuildsGetAll godoc
//
//	@Summary	Get All Builds
//	@Tags		Builds
//	@Produce	json
//	@Success	200	{object}	Results[repo.BuildOut]
//	@Router		/v1/builds [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildsGetAll() errchain.HandlerFunc {
	fn := func(r *http.Request, _ struct{}) ([]repo.BuildOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.GetAll(auth, auth.GID)
	}

	return adapters.Query(fn, http.StatusOK)
}

// HandleBuildsCreate godoc
//
//	@Summary	Create Build
//	@Tags		Builds
//	@Produce	json
//	@Param		payload	body	repo.BuildCreate	true	"Build Data"
//	@Success	201		{object}	repo.BuildOut
//	@Router		/v1/builds [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildsCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, body repo.BuildCreate) (repo.BuildOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.Create(auth, auth.GID, body)
	}

	return adapters.Action(fn, http.StatusCreated)
}

// HandleBuildGet godoc
//
//	@Summary	Get Build
//	@Tags		Builds
//	@Produce	json
//	@Param		id	path	string	true	"Build ID"
//	@Success	200	{object}	repo.BuildOut
//	@Router		/v1/builds/{id} [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) (repo.BuildOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.GetOne(auth, auth.GID, ID)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

// HandleBuildUpdate godoc
//
//	@Summary	Update Build
//	@Tags		Builds
//	@Produce	json
//	@Param		id		path	string				true	"Build ID"
//	@Param		payload	body	repo.BuildUpdate	true	"Build Data"
//	@Success	200		{object}	repo.BuildOut
//	@Router		/v1/builds/{id} [PUT]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildUpdate() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, body repo.BuildUpdate) (repo.BuildOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.Update(auth, auth.GID, ID, body)
	}

	return adapters.ActionID("id", fn, http.StatusOK)
}

// HandleBuildDelete godoc
//
//	@Summary	Delete Build
//	@Tags		Builds
//	@Param		id	path	string	true	"Build ID"
//	@Success	204
//	@Router		/v1/builds/{id} [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildDelete() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID) (any, error) {
		auth := services.NewContext(r.Context())
		return nil, ctrl.repo.Builds.Delete(auth, auth.GID, ID)
	}

	return adapters.CommandID("id", fn, http.StatusNoContent)
}

// HandleBuildPartsGet godoc
//
//	@Summary	Get Build Parts (mounting history)
//	@Tags		Builds
//	@Produce	json
//	@Param		id	path	string	true	"Build ID"
//	@Success	200	{object}	Results[repo.BuildPartOut]
//	@Router		/v1/builds/{id}/parts [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildPartsGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) ([]repo.BuildPartOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.GetParts(auth, auth.GID, ID, false)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

// HandleBuildPartsCurrentGet godoc
//
//	@Summary	Get Currently Mounted Parts
//	@Tags		Builds
//	@Produce	json
//	@Param		id	path	string	true	"Build ID"
//	@Success	200	{object}	Results[repo.BuildPartOut]
//	@Router		/v1/builds/{id}/parts/current [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildPartsCurrentGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) ([]repo.BuildPartOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.GetParts(auth, auth.GID, ID, true)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

// HandleBuildPartAdd godoc
//
//	@Summary	Mount a Part on a Build
//	@Tags		Builds
//	@Produce	json
//	@Param		id		path	string					true	"Build ID"
//	@Param		payload	body	repo.BuildPartCreate	true	"Part Data"
//	@Success	201		{object}	repo.BuildPartAddOut
//	@Router		/v1/builds/{id}/parts [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildPartAdd() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, body repo.BuildPartCreate) (repo.BuildPartAddOut, error) {
		auth := services.NewContext(r.Context())
		out, err := ctrl.repo.Builds.AddPart(auth, auth.GID, ID, body)
		if errors.Is(err, repo.ErrBuildRetired) {
			return repo.BuildPartAddOut{}, validate.NewRequestError(err, http.StatusConflict)
		}
		return out, err
	}

	return adapters.ActionID("id", fn, http.StatusCreated)
}

// HandleBuildSports godoc
//
//	@Summary	Get Supported Build Sports
//	@Tags		Builds
//	@Produce	json
//	@Success	200	{array}	repo.SportPositionConfig
//	@Router		/v1/builds/sports [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildSports() errchain.HandlerFunc {
	fn := func(r *http.Request, _ struct{}) ([]repo.SportPositionConfig, error) {
		return ctrl.repo.Builds.Sports(), nil
	}

	return adapters.Query(fn, http.StatusOK)
}

// HandleBuildsStats godoc
//
//	@Summary	Get Build Statistics
//	@Tags		Builds
//	@Produce	json
//	@Success	200	{object}	repo.BuildsStatsOut
//	@Router		/v1/builds/statistics [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildsStats() errchain.HandlerFunc {
	fn := func(r *http.Request, _ struct{}) (repo.BuildsStatsOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.Stats(auth, auth.GID)
	}

	return adapters.Query(fn, http.StatusOK)
}

// HandleBuildPartRemove godoc
//
//	@Summary	Remove a Part from a Build (stamps removedDate)
//	@Tags		Builds
//	@Param		id			path	string	true	"Build ID"
//	@Param		part_id		path	string	true	"Build Part ID"
//	@Success	204
//	@Router		/v1/builds/{id}/parts/{part_id} [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildPartRemove() errchain.HandlerFunc {
	fn := func(r *http.Request, partID uuid.UUID) (any, error) {
		auth := services.NewContext(r.Context())
		return nil, ctrl.repo.Builds.RemovePart(auth, auth.GID, partID)
	}

	return adapters.CommandID("part_id", fn, http.StatusNoContent)
}

// HandleBuildPartDelete godoc
//
//	@Summary	Permanently delete a Build Part assignment
//	@Tags		Builds
//	@Param		id			path	string	true	"Build ID"
//	@Param		part_id		path	string	true	"Build Part ID"
//	@Success	204
//	@Router		/v1/builds/{id}/parts/{part_id}/purge [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildPartDelete() errchain.HandlerFunc {
	fn := func(r *http.Request, partID uuid.UUID) (any, error) {
		auth := services.NewContext(r.Context())
		return nil, ctrl.repo.Builds.DeletePart(auth, auth.GID, partID)
	}

	return adapters.CommandID("part_id", fn, http.StatusNoContent)
}

// HandleBuildSessionsGet godoc
//
//	@Summary	Get Usage Sessions for a Build
//	@Tags		Builds
//	@Produce	json
//	@Param		id	path	string	true	"Build ID"
//	@Success	200	{object}	Results[repo.UsageSessionOut]
//	@Router		/v1/builds/{id}/sessions [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildSessionsGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) ([]repo.UsageSessionOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.GetSessions(auth, auth.GID, ID)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

// HandleBuildSessionCreate godoc
//
//	@Summary	Log a Usage Session
//	@Tags		Builds
//	@Produce	json
//	@Param		id		path	string					true	"Build ID"
//	@Param		payload	body	repo.UsageSessionCreate	true	"Session Data"
//	@Success	201		{object}	repo.UsageSessionOut
//	@Router		/v1/builds/{id}/sessions [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildSessionCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, body repo.UsageSessionCreate) (repo.UsageSessionOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.CreateSession(auth, auth.GID, ID, body)
	}

	return adapters.ActionID("id", fn, http.StatusCreated)
}

// HandleBuildSessionUpdate godoc
//
//	@Summary	Update a Usage Session
//	@Tags		Builds
//	@Produce	json
//	@Param		id			path	string					true	"Build ID"
//	@Param		session_id	path	string					true	"Session ID"
//	@Param		payload		body	repo.UsageSessionUpdate	true	"Session Data"
//	@Success	200		{object}	repo.UsageSessionOut
//	@Router		/v1/builds/{id}/sessions/{session_id} [PUT]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildSessionUpdate() errchain.HandlerFunc {
	fn := func(r *http.Request, sessionID uuid.UUID, body repo.UsageSessionUpdate) (repo.UsageSessionOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.UpdateSession(auth, auth.GID, sessionID, body)
	}

	return adapters.ActionID("session_id", fn, http.StatusOK)
}

// HandleBuildSessionDelete godoc
//
//	@Summary	Delete a Usage Session
//	@Tags		Builds
//	@Param		id			path	string	true	"Build ID"
//	@Param		session_id	path	string	true	"Session ID"
//	@Success	204
//	@Router		/v1/builds/{id}/sessions/{session_id} [DELETE]
//	@Security	Bearer
func (ctrl *V1Controller) HandleBuildSessionDelete() errchain.HandlerFunc {
	fn := func(r *http.Request, sessionID uuid.UUID) (any, error) {
		auth := services.NewContext(r.Context())
		return nil, ctrl.repo.Builds.DeleteSession(auth, auth.GID, sessionID)
	}

	return adapters.CommandID("session_id", fn, http.StatusNoContent)
}

// HandleEntityBuildHistoryGet godoc
//
//	@Summary	Get the mounting history of a part across all builds
//	@Tags		Builds
//	@Produce	json
//	@Param		id	path	string	true	"Entity (Part) ID"
//	@Success	200	{object}	Results[repo.BuildPartOut]
//	@Router		/v1/entities/{id}/build-history [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandleEntityBuildHistoryGet() errchain.HandlerFunc {
	fn := func(r *http.Request, ID uuid.UUID, _ struct{}) ([]repo.BuildPartOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Builds.GetPartsByEntity(auth, auth.GID, ID)
	}

	return adapters.QueryID("id", fn, http.StatusOK)
}

package repo

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/build"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/buildpart"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/entity"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/group"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/predicate"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/usagesession"
)

// BuildRepository covers builds, their part assignments, and usage sessions.
// A BuildPart row with a null removed_date means the part is currently
// mounted; stamped rows preserve the mounting history across builds.
type BuildRepository struct {
	db *ent.Client
}

type BuildCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Sport       string `json:"sport"`
	Status      string `json:"status" validate:"omitempty,oneof=assembled disassembled retired"`
}

type BuildUpdate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Sport       string `json:"sport"`
	Status      string `json:"status" validate:"omitempty,oneof=assembled disassembled retired"`
}

type BuildOut struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Sport       string    `json:"sport"`
	Status      string    `json:"status"`
}

type BuildPartCreate struct {
	// EntityID is the part (an Entity) being mounted.
	EntityID uuid.UUID `json:"entityId" validate:"required"`
	// Position is required; a mounted part must live at a concrete position.
	Position string `json:"position" validate:"required"`
}

type BuildPartOut struct {
	ID          uuid.UUID  `json:"id"`
	BuildID     uuid.UUID  `json:"buildId"`
	EntityID    uuid.UUID  `json:"entityId"`
	EntityName  string     `json:"entityName"`
	Position    string     `json:"position"`
	AddedDate   time.Time  `json:"addedDate"`
	RemovedDate *time.Time `json:"removedDate"`
}

// BuildPartAddOut wraps a freshly mounted part with non-blocking warnings so
// the client can surface them without failing the mount.
type BuildPartAddOut struct {
	Part     BuildPartOut `json:"part"`
	Warnings []string     `json:"warnings"`
}

// BuildsStatsOut aggregates investment and usage analytics across all of the
// group's builds. CostPerSession is 0 when no sessions have been logged.
type BuildsStatsOut struct {
	TotalInvestment float64 `json:"totalInvestment"`
	SessionCount    int     `json:"sessionCount"`
	CostPerSession  float64 `json:"costPerSession"`
}

type UsageSessionCreate struct {
	Date            time.Time `json:"date"`
	DurationMinutes int       `json:"durationMinutes"`
	Location        string    `json:"location"`
	Conditions      string    `json:"conditions"`
	Notes           string    `json:"notes"`
	Rating          int       `json:"rating"`
}

type UsageSessionUpdate struct {
	Date            time.Time `json:"date"`
	DurationMinutes int       `json:"durationMinutes"`
	Location        string    `json:"location"`
	Conditions      string    `json:"conditions"`
	Notes           string    `json:"notes"`
	Rating          int       `json:"rating"`
}

type UsageSessionOut struct {
	ID              uuid.UUID `json:"id"`
	BuildID         uuid.UUID `json:"buildId"`
	Date            time.Time `json:"date"`
	DurationMinutes int       `json:"durationMinutes"`
	Location        string    `json:"location"`
	Conditions      string    `json:"conditions"`
	Notes           string    `json:"notes"`
	Rating          int       `json:"rating"`
}

// ErrBuildRetired is returned when trying to mount a part on a retired build.
var ErrBuildRetired = errors.New("build is retired; set its status to assembled or disassembled to mount parts")

var (
	mapBuildErr  = mapTErrFunc(mapBuild)
	mapEachBuild = mapTEachFunc(mapBuild)
)

func mapBuild(s *ent.Build) BuildOut {
	return BuildOut{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Sport:       s.Sport,
		Status:      string(s.Status),
	}
}

func mapBuildPart(p *ent.BuildPart) BuildPartOut {
	out := BuildPartOut{
		ID:          p.ID,
		BuildID:     p.BuildID,
		EntityID:    p.EntityID,
		Position:    p.Position,
		AddedDate:   p.AddedDate,
		RemovedDate: nil,
	}
	if !p.RemovedDate.IsZero() {
		removed := p.RemovedDate
		out.RemovedDate = &removed
	}
	if p.Edges.Entity != nil {
		out.EntityName = p.Edges.Entity.Name
	}
	return out
}

func mapUsageSession(s *ent.UsageSession) UsageSessionOut {
	return UsageSessionOut{
		ID:              s.ID,
		BuildID:         s.BuildID,
		Date:            s.Date,
		DurationMinutes: s.DurationMinutes,
		Location:        s.Location,
		Conditions:      s.Conditions,
		Notes:           s.Notes,
		Rating:          s.Rating,
	}
}

// buildScoped returns a predicate limiting builds to the given group.
func buildScoped(gid uuid.UUID) predicate.Build {
	return build.HasGroupWith(group.ID(gid))
}

func (r *BuildRepository) GetAll(ctx context.Context, gid uuid.UUID) ([]BuildOut, error) {
	builds, err := r.db.Build.Query().
		Where(buildScoped(gid)).
		Order(ent.Desc(build.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return mapEachBuild(builds), nil
}

func (r *BuildRepository) GetOne(ctx context.Context, gid, id uuid.UUID) (BuildOut, error) {
	s, err := r.db.Build.Query().
		Where(build.ID(id), buildScoped(gid)).
		Only(ctx)
	if err != nil {
		return BuildOut{}, err
	}
	return mapBuild(s), nil
}

func (r *BuildRepository) Create(ctx context.Context, gid uuid.UUID, input BuildCreate) (BuildOut, error) {
	status := build.Status(input.Status)
	if input.Status == "" {
		status = build.StatusAssembled
	}
	s, err := r.db.Build.Create().
		SetGroupID(gid).
		SetName(input.Name).
		SetDescription(input.Description).
		SetSport(input.Sport).
		SetStatus(status).
		Save(ctx)
	return mapBuildErr(s, err)
}

func (r *BuildRepository) Update(ctx context.Context, gid, id uuid.UUID, input BuildUpdate) (BuildOut, error) {
	update := r.db.Build.Update().
		Where(build.ID(id), buildScoped(gid)).
		SetName(input.Name).
		SetDescription(input.Description).
		SetSport(input.Sport)
	if input.Status != "" {
		update.SetStatus(build.Status(input.Status))
	}
	s, err := update.Save(ctx)
	if err != nil {
		return BuildOut{}, err
	}
	if s == 0 {
		return BuildOut{}, &ent.NotFoundError{}
	}
	return r.GetOne(ctx, gid, id)
}

func (r *BuildRepository) Delete(ctx context.Context, gid, id uuid.UUID) error {
	deleted, err := r.db.Build.Delete().
		Where(build.ID(id), buildScoped(gid)).
		Exec(ctx)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return &ent.NotFoundError{}
	}
	return nil
}

// AddPart mounts a part on the build. The part (Entity) must belong to the
// caller's group, mirroring the ownership check in MaintenanceEntryRepository.
// Mounting on a retired build is rejected (ErrBuildRetired). Non-blocking
// warnings (machine-readable codes) are returned alongside the created
// assignment; they never fail the mount.
func (r *BuildRepository) AddPart(ctx context.Context, gid, buildID uuid.UUID, input BuildPartCreate) (BuildPartAddOut, error) {
	e, err := r.db.Entity.Query().
		Where(entity.ID(input.EntityID), entity.HasGroupWith(group.ID(gid))).
		WithEntityType().
		Only(ctx)
	if err != nil {
		return BuildPartAddOut{}, err
	}

	b, err := r.db.Build.Query().
		Where(build.ID(buildID), buildScoped(gid)).
		Only(ctx)
	if err != nil {
		return BuildPartAddOut{}, err
	}
	if b.Status == build.StatusRetired {
		return BuildPartAddOut{}, ErrBuildRetired
	}

	warnings := []string{}
	if et := e.Edges.EntityType; et != nil && et.IsLocation {
		warnings = append(warnings, "entity_is_location")
	}

	// A part can only be physically mounted on one build at a time; warn
	// (but don't block) when it is currently mounted on another build.
	current, err := r.db.BuildPart.Query().
		Where(
			buildpart.EntityID(input.EntityID),
			buildpart.RemovedDateIsNil(),
			buildpart.HasBuildWith(buildScoped(gid)),
		).
		All(ctx)
	if err != nil {
		return BuildPartAddOut{}, err
	}
	for _, cp := range current {
		if cp.BuildID != buildID {
			warnings = append(warnings, "part_mounted_elsewhere")
			break
		}
	}

	if b.Status == build.StatusDisassembled {
		warnings = append(warnings, "build_disassembled")
	}

	// Warn (but don't block) when a part made for one sport is mounted on a
	// build for another. Unassigned (empty) and explicit "generic" parts are
	// never flagged.
	if b.Sport != "" && e.Sport != "" && e.Sport != "generic" && b.Sport != e.Sport {
		warnings = append(warnings, "part_sport_mismatch")
	}

	// Soft-validate the position against the build's sport; unknown sports
	// or empty positions are left untouched.
	if input.Position != "" {
		if positions := SportPositions(b.Sport); positions != nil && !slices.Contains(positions, input.Position) {
			warnings = append(warnings, "unknown_position")
		}
	}

	created, err := r.db.BuildPart.Create().
		SetBuildID(buildID).
		SetEntityID(input.EntityID).
		SetNillablePosition(&input.Position).
		Save(ctx)
	if err != nil {
		return BuildPartAddOut{}, err
	}

	part, err := r.GetPart(ctx, gid, created.ID)
	if err != nil {
		return BuildPartAddOut{}, err
	}

	return BuildPartAddOut{Part: part, Warnings: warnings}, nil
}

func (r *BuildRepository) GetPart(ctx context.Context, gid, id uuid.UUID) (BuildPartOut, error) {
	p, err := r.db.BuildPart.Query().
		Where(
			buildpart.ID(id),
			buildpart.HasBuildWith(buildScoped(gid)),
		).
		WithEntity().
		Only(ctx)
	if err != nil {
		return BuildPartOut{}, err
	}
	return mapBuildPart(p), nil
}

// GetParts returns the part assignments of a build. When currentOnly is true,
// only currently mounted parts (removed_date IS NULL) are returned.
func (r *BuildRepository) GetParts(ctx context.Context, gid, buildID uuid.UUID, currentOnly bool) ([]BuildPartOut, error) {
	query := r.db.BuildPart.Query().
		Where(
			buildpart.BuildID(buildID),
			buildpart.HasBuildWith(buildScoped(gid)),
		).
		WithEntity().
		Order(ent.Asc(buildpart.FieldAddedDate))
	if currentOnly {
		query = query.Where(buildpart.RemovedDateIsNil())
	}
	parts, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]BuildPartOut, len(parts))
	for i, p := range parts {
		out[i] = mapBuildPart(p)
	}
	return out, nil
}

// RemovePart stamps removed_date on the assignment rather than deleting the
// row, so the mounting history survives. Deleting the row outright would
// silently erase evidence that the part was ever in the build.
func (r *BuildRepository) RemovePart(ctx context.Context, gid, id uuid.UUID) error {
	updated, err := r.db.BuildPart.Update().
		Where(
			buildpart.ID(id),
			buildpart.RemovedDateIsNil(),
			buildpart.HasBuildWith(buildScoped(gid)),
		).
		SetRemovedDate(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	if updated == 0 {
		return &ent.NotFoundError{}
	}
	return nil
}

// DeletePart permanently removes a build-part assignment row.
func (r *BuildRepository) DeletePart(ctx context.Context, gid, id uuid.UUID) error {
	deleted, err := r.db.BuildPart.Delete().
		Where(
			buildpart.ID(id),
			buildpart.HasBuildWith(buildScoped(gid)),
		).
		Exec(ctx)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return &ent.NotFoundError{}
	}
	return nil
}

// GetPartsByEntity returns the full mounting history of one part across all
// of the group's builds, newest first.
func (r *BuildRepository) GetPartsByEntity(ctx context.Context, gid, entityID uuid.UUID) ([]BuildPartOut, error) {
	parts, err := r.db.BuildPart.Query().
		Where(
			buildpart.EntityID(entityID),
			buildpart.HasBuildWith(buildScoped(gid)),
		).
		WithEntity().
		Order(ent.Desc(buildpart.FieldAddedDate)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]BuildPartOut, len(parts))
	for i, p := range parts {
		out[i] = mapBuildPart(p)
	}
	return out, nil
}

var (
	mapUsageSessionErr  = mapTErrFunc(mapUsageSession)
	mapEachUsageSession = mapTEachFunc(mapUsageSession)
)

// Stats aggregates analytics across all of the group's builds: total
// investment in currently mounted parts (sum of distinct part purchase
// prices), logged session count, and cost per session (0 when no sessions).
func (r *BuildRepository) Stats(ctx context.Context, gid uuid.UUID) (BuildsStatsOut, error) {
	parts, err := r.db.BuildPart.Query().
		Where(
			buildpart.RemovedDateIsNil(),
			buildpart.HasBuildWith(buildScoped(gid)),
		).
		WithEntity().
		All(ctx)
	if err != nil {
		return BuildsStatsOut{}, err
	}

	// A part can be mounted in several builds at once; count its purchase
	// price only once toward the investment total.
	seen := make(map[uuid.UUID]struct{}, len(parts))
	var investment float64
	for _, p := range parts {
		e := p.Edges.Entity
		if e == nil {
			continue
		}
		if _, ok := seen[e.ID]; ok {
			continue
		}
		seen[e.ID] = struct{}{}
		investment += e.PurchasePrice
	}

	sessionCount, err := r.db.UsageSession.Query().
		Where(usagesession.HasBuildWith(buildScoped(gid))).
		Count(ctx)
	if err != nil {
		return BuildsStatsOut{}, err
	}

	costPerSession := 0.0
	if sessionCount > 0 {
		costPerSession = investment / float64(sessionCount)
	}

	return BuildsStatsOut{
		TotalInvestment: investment,
		SessionCount:    sessionCount,
		CostPerSession:  costPerSession,
	}, nil
}

func (r *BuildRepository) GetSessions(ctx context.Context, gid, buildID uuid.UUID) ([]UsageSessionOut, error) {
	sessions, err := r.db.UsageSession.Query().
		Where(
			usagesession.BuildID(buildID),
			usagesession.HasBuildWith(buildScoped(gid)),
		).
		Order(ent.Desc(usagesession.FieldDate), ent.Desc(usagesession.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return mapEachUsageSession(sessions), nil
}

func (r *BuildRepository) CreateSession(ctx context.Context, gid, buildID uuid.UUID, input UsageSessionCreate) (UsageSessionOut, error) {
	created, err := r.db.UsageSession.Create().
		SetBuildID(buildID).
		SetDate(input.Date).
		SetDurationMinutes(input.DurationMinutes).
		SetNillableLocation(&input.Location).
		SetNillableConditions(&input.Conditions).
		SetNillableNotes(&input.Notes).
		SetRating(input.Rating).
		Save(ctx)
	return mapUsageSessionErr(created, err)
}

func (r *BuildRepository) UpdateSession(ctx context.Context, gid, id uuid.UUID, input UsageSessionUpdate) (UsageSessionOut, error) {
	updated, err := r.db.UsageSession.Update().
		Where(
			usagesession.ID(id),
			usagesession.HasBuildWith(buildScoped(gid)),
		).
		SetDate(input.Date).
		SetDurationMinutes(input.DurationMinutes).
		SetNillableLocation(&input.Location).
		SetNillableConditions(&input.Conditions).
		SetNillableNotes(&input.Notes).
		SetRating(input.Rating).
		Save(ctx)
	if err != nil {
		return UsageSessionOut{}, err
	}
	if updated == 0 {
		return UsageSessionOut{}, &ent.NotFoundError{}
	}
	session, err := r.db.UsageSession.Get(ctx, id)
	if err != nil {
		return UsageSessionOut{}, err
	}
	return mapUsageSession(session), nil
}

func (r *BuildRepository) DeleteSession(ctx context.Context, gid, id uuid.UUID) error {
	deleted, err := r.db.UsageSession.Delete().
		Where(
			usagesession.ID(id),
			usagesession.HasBuildWith(buildScoped(gid)),
		).
		Exec(ctx)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return &ent.NotFoundError{}
	}
	return nil
}

package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/entity"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/group"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/predicate"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/setup"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/setuppart"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent/usagesession"
)

// SetupRepository covers setups, their part assignments, and usage sessions.
// A SetupPart row with a null removed_date means the part is currently
// mounted; stamped rows preserve the mounting history across setups.
type SetupRepository struct {
	db *ent.Client
}

type SetupCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Sport       string `json:"sport"`
	IsActive    bool   `json:"isActive"`
}

type SetupUpdate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Sport       string `json:"sport"`
	IsActive    bool   `json:"isActive"`
}

type SetupOut struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Sport       string    `json:"sport"`
	IsActive    bool      `json:"isActive"`
}

type SetupPartCreate struct {
	// EntityID is the part (an Entity) being mounted.
	EntityID uuid.UUID `json:"entityId" validate:"required"`
	Position string    `json:"position"`
}

type SetupPartOut struct {
	ID          uuid.UUID  `json:"id"`
	SetupID     uuid.UUID  `json:"setupId"`
	EntityID    uuid.UUID  `json:"entityId"`
	EntityName  string     `json:"entityName"`
	Position    string     `json:"position"`
	AddedDate   time.Time  `json:"addedDate"`
	RemovedDate *time.Time `json:"removedDate"`
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
	SetupID         uuid.UUID `json:"setupId"`
	Date            time.Time `json:"date"`
	DurationMinutes int       `json:"durationMinutes"`
	Location        string    `json:"location"`
	Conditions      string    `json:"conditions"`
	Notes           string    `json:"notes"`
	Rating          int       `json:"rating"`
}

var (
	mapSetupErr  = mapTErrFunc(mapSetup)
	mapEachSetup = mapTEachFunc(mapSetup)
)

func mapSetup(s *ent.Setup) SetupOut {
	return SetupOut{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Sport:       s.Sport,
		IsActive:    s.IsActive,
	}
}

func mapSetupPart(p *ent.SetupPart) SetupPartOut {
	out := SetupPartOut{
		ID:          p.ID,
		SetupID:     p.SetupID,
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
		SetupID:         s.SetupID,
		Date:            s.Date,
		DurationMinutes: s.DurationMinutes,
		Location:        s.Location,
		Conditions:      s.Conditions,
		Notes:           s.Notes,
		Rating:          s.Rating,
	}
}

// setupScoped returns a predicate limiting setups to the given group.
func setupScoped(gid uuid.UUID) predicate.Setup {
	return setup.HasGroupWith(group.ID(gid))
}

func (r *SetupRepository) GetAll(ctx context.Context, gid uuid.UUID) ([]SetupOut, error) {
	setups, err := r.db.Setup.Query().
		Where(setupScoped(gid)).
		Order(ent.Desc(setup.FieldIsActive), ent.Desc(setup.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return mapEachSetup(setups), nil
}

func (r *SetupRepository) GetOne(ctx context.Context, gid, id uuid.UUID) (SetupOut, error) {
	s, err := r.db.Setup.Query().
		Where(setup.ID(id), setupScoped(gid)).
		Only(ctx)
	if err != nil {
		return SetupOut{}, err
	}
	return mapSetup(s), nil
}

func (r *SetupRepository) Create(ctx context.Context, gid uuid.UUID, input SetupCreate) (SetupOut, error) {
	s, err := r.db.Setup.Create().
		SetGroupID(gid).
		SetName(input.Name).
		SetDescription(input.Description).
		SetSport(input.Sport).
		SetIsActive(input.IsActive).
		Save(ctx)
	return mapSetupErr(s, err)
}

func (r *SetupRepository) Update(ctx context.Context, gid, id uuid.UUID, input SetupUpdate) (SetupOut, error) {
	s, err := r.db.Setup.Update().
		Where(setup.ID(id), setupScoped(gid)).
		SetName(input.Name).
		SetDescription(input.Description).
		SetSport(input.Sport).
		SetIsActive(input.IsActive).
		Save(ctx)
	if err != nil {
		return SetupOut{}, err
	}
	if s == 0 {
		return SetupOut{}, &ent.NotFoundError{}
	}
	return r.GetOne(ctx, gid, id)
}

func (r *SetupRepository) Delete(ctx context.Context, gid, id uuid.UUID) error {
	deleted, err := r.db.Setup.Delete().
		Where(setup.ID(id), setupScoped(gid)).
		Exec(ctx)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return &ent.NotFoundError{}
	}
	return nil
}

// AddPart mounts a part on the setup. The part (Entity) must belong to the
// caller's group, mirroring the ownership check in MaintenanceEntryRepository.
func (r *SetupRepository) AddPart(ctx context.Context, gid, setupID uuid.UUID, input SetupPartCreate) (SetupPartOut, error) {
	owned, err := r.db.Entity.Query().
		Where(entity.ID(input.EntityID), entity.HasGroupWith(group.ID(gid))).
		Exist(ctx)
	if err != nil {
		return SetupPartOut{}, err
	}
	if !owned {
		return SetupPartOut{}, &ent.NotFoundError{}
	}

	created, err := r.db.SetupPart.Create().
		SetSetupID(setupID).
		SetEntityID(input.EntityID).
		SetNillablePosition(&input.Position).
		Save(ctx)
	if err != nil {
		return SetupPartOut{}, err
	}

	return r.GetPart(ctx, gid, created.ID)
}

func (r *SetupRepository) GetPart(ctx context.Context, gid, id uuid.UUID) (SetupPartOut, error) {
	p, err := r.db.SetupPart.Query().
		Where(
			setuppart.ID(id),
			setuppart.HasSetupWith(setupScoped(gid)),
		).
		WithEntity().
		Only(ctx)
	if err != nil {
		return SetupPartOut{}, err
	}
	return mapSetupPart(p), nil
}

// GetParts returns the part assignments of a setup. When currentOnly is true,
// only currently mounted parts (removed_date IS NULL) are returned.
func (r *SetupRepository) GetParts(ctx context.Context, gid, setupID uuid.UUID, currentOnly bool) ([]SetupPartOut, error) {
	query := r.db.SetupPart.Query().
		Where(
			setuppart.SetupID(setupID),
			setuppart.HasSetupWith(setupScoped(gid)),
		).
		WithEntity().
		Order(ent.Asc(setuppart.FieldAddedDate))
	if currentOnly {
		query = query.Where(setuppart.RemovedDateIsNil())
	}
	parts, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]SetupPartOut, len(parts))
	for i, p := range parts {
		out[i] = mapSetupPart(p)
	}
	return out, nil
}

// RemovePart stamps removed_date on the assignment rather than deleting the
// row, so the mounting history survives. Deleting the row outright would
// silently erase evidence that the part was ever in the setup.
func (r *SetupRepository) RemovePart(ctx context.Context, gid, id uuid.UUID) error {
	updated, err := r.db.SetupPart.Update().
		Where(
			setuppart.ID(id),
			setuppart.RemovedDateIsNil(),
			setuppart.HasSetupWith(setupScoped(gid)),
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

// DeletePart permanently removes a setup-part assignment row.
func (r *SetupRepository) DeletePart(ctx context.Context, gid, id uuid.UUID) error {
	deleted, err := r.db.SetupPart.Delete().
		Where(
			setuppart.ID(id),
			setuppart.HasSetupWith(setupScoped(gid)),
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
// of the group's setups, newest first.
func (r *SetupRepository) GetPartsByEntity(ctx context.Context, gid, entityID uuid.UUID) ([]SetupPartOut, error) {
	parts, err := r.db.SetupPart.Query().
		Where(
			setuppart.EntityID(entityID),
			setuppart.HasSetupWith(setupScoped(gid)),
		).
		WithEntity().
		Order(ent.Desc(setuppart.FieldAddedDate)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]SetupPartOut, len(parts))
	for i, p := range parts {
		out[i] = mapSetupPart(p)
	}
	return out, nil
}

var (
	mapUsageSessionErr  = mapTErrFunc(mapUsageSession)
	mapEachUsageSession = mapTEachFunc(mapUsageSession)
)

func (r *SetupRepository) GetSessions(ctx context.Context, gid, setupID uuid.UUID) ([]UsageSessionOut, error) {
	sessions, err := r.db.UsageSession.Query().
		Where(
			usagesession.SetupID(setupID),
			usagesession.HasSetupWith(setupScoped(gid)),
		).
		Order(ent.Desc(usagesession.FieldDate), ent.Desc(usagesession.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return mapEachUsageSession(sessions), nil
}

func (r *SetupRepository) CreateSession(ctx context.Context, gid, setupID uuid.UUID, input UsageSessionCreate) (UsageSessionOut, error) {
	created, err := r.db.UsageSession.Create().
		SetSetupID(setupID).
		SetDate(input.Date).
		SetDurationMinutes(input.DurationMinutes).
		SetNillableLocation(&input.Location).
		SetNillableConditions(&input.Conditions).
		SetNillableNotes(&input.Notes).
		SetRating(input.Rating).
		Save(ctx)
	return mapUsageSessionErr(created, err)
}

func (r *SetupRepository) UpdateSession(ctx context.Context, gid, id uuid.UUID, input UsageSessionUpdate) (UsageSessionOut, error) {
	updated, err := r.db.UsageSession.Update().
		Where(
			usagesession.ID(id),
			usagesession.HasSetupWith(setupScoped(gid)),
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

func (r *SetupRepository) DeleteSession(ctx context.Context, gid, id uuid.UUID) error {
	deleted, err := r.db.UsageSession.Delete().
		Where(
			usagesession.ID(id),
			usagesession.HasSetupWith(setupScoped(gid)),
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

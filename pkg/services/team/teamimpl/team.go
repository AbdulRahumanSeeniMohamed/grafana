package teamimpl

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/services/team"
	"github.com/grafana/grafana/pkg/setting"
)

type Service struct {
	store  store
	tracer tracing.Tracer
}

func ProvideService(db db.ReplDB, cfg *setting.Cfg, tracer tracing.Tracer) (team.Service, error) {
	return &Service{
		store:  &xormStore{db: db, cfg: cfg, deletes: []string{}},
		tracer: tracer,
	}, nil
}

func (s *Service) CreateTeam(ctx context.Context, name, email string, orgID int64) (team.Team, error) {
	_, span := s.tracer.Start(ctx, "team.CreateTeam", trace.WithAttributes(
		attribute.Int64("orgID", orgID),
		attribute.String("name", name),
	))
	defer span.End()
	return s.store.Create(name, email, orgID)
}

func (s *Service) UpdateTeam(ctx context.Context, cmd *team.UpdateTeamCommand) error {
	ctx, span := s.tracer.Start(ctx, "team.UpdateTeam", trace.WithAttributes(
		attribute.Int64("orgID", cmd.OrgID),
		attribute.Int64("teamID", cmd.ID),
	))
	defer span.End()
	return s.store.Update(ctx, cmd)
}

func (s *Service) DeleteTeam(ctx context.Context, cmd *team.DeleteTeamCommand) error {
	ctx, span := s.tracer.Start(ctx, "team.DeleteTeam", trace.WithAttributes(
		attribute.Int64("orgID", cmd.OrgID),
		attribute.Int64("teamID", cmd.ID),
	))
	defer span.End()
	return s.store.Delete(ctx, cmd)
}

func (s *Service) SearchTeams(ctx context.Context, query *team.SearchTeamsQuery) (team.SearchTeamQueryResult, error) {
	ctx, span := s.tracer.Start(ctx, "team.SearchTeams", trace.WithAttributes(
		attribute.Int64("orgID", query.OrgID),
		attribute.String("query", query.Query),
	))
	defer span.End()
	return s.store.Search(ctx, query)
}

func (s *Service) GetTeamByID(ctx context.Context, query *team.GetTeamByIDQuery) (*team.TeamDTO, error) {
	ctx, span := s.tracer.Start(ctx, "team.GetTeamByID", trace.WithAttributes(
		attribute.Int64("orgID", query.OrgID),
		attribute.Int64("teamID", query.ID),
	))
	defer span.End()
	return s.store.GetByID(ctx, query)
}

func (s *Service) GetTeamsByUser(ctx context.Context, query *team.GetTeamsByUserQuery) ([]*team.TeamDTO, error) {
	ctx, span := s.tracer.Start(ctx, "team.GetTeamsByUser", trace.WithAttributes(
		attribute.Int64("orgID", query.OrgID),
		attribute.Int64("userID", query.UserID),
	))
	defer span.End()
	return s.store.GetByUser(ctx, query)
}

func (s *Service) GetTeamIDsByUser(ctx context.Context, query *team.GetTeamIDsByUserQuery) ([]int64, error) {
	ctx, span := s.tracer.Start(ctx, "team.GetTeamIDsByUser", trace.WithAttributes(
		attribute.Int64("orgID", query.OrgID),
		attribute.Int64("userID", query.UserID),
	))
	defer span.End()
	return s.store.GetIDsByUser(ctx, query)
}

func (s *Service) IsTeamMember(ctx context.Context, orgId int64, teamId int64, userId int64) (bool, error) {
	_, span := s.tracer.Start(ctx, "team.IsTeamMember", trace.WithAttributes(
		attribute.Int64("orgID", orgId),
		attribute.Int64("teamID", teamId),
		attribute.Int64("userID", userId),
	))
	defer span.End()
	return s.store.IsMember(orgId, teamId, userId)
}

func (s *Service) RemoveUsersMemberships(ctx context.Context, userID int64) error {
	ctx, span := s.tracer.Start(ctx, "team.RemoveUsersMemberships", trace.WithAttributes(
		attribute.Int64("userID", userID),
	))
	defer span.End()
	return s.store.RemoveUsersMemberships(ctx, userID)
}

func (s *Service) GetUserTeamMemberships(ctx context.Context, orgID, userID int64, external bool) ([]*team.TeamMemberDTO, error) {
	ctx, span := s.tracer.Start(ctx, "team.GetUserTeamMemberships", trace.WithAttributes(
		attribute.Int64("orgID", orgID),
		attribute.Int64("userID", userID),
	))
	defer span.End()
	return s.store.GetMemberships(ctx, orgID, userID, external)
}

func (s *Service) GetTeamMembers(ctx context.Context, query *team.GetTeamMembersQuery) ([]*team.TeamMemberDTO, error) {
	ctx, span := s.tracer.Start(ctx, "team.GetTeamMembers", trace.WithAttributes(
		attribute.Int64("orgID", query.OrgID),
		attribute.Int64("teamID", query.TeamID),
		attribute.String("teamUID", query.TeamUID),
	))
	defer span.End()
	return s.store.GetMembers(ctx, query)
}

func (s *Service) RegisterDelete(query string) {
	s.store.RegisterDelete(query)
}

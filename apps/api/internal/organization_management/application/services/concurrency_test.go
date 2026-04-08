package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/app_errors"
	"league-management/internal/shared/dtos"
	user_domain "league-management/internal/user_management/domain"
	"sync"
	"testing"
)

type fakeUserRepo struct {
	users map[string]*user_domain.User
}

func (r *fakeUserRepo) FindById(id string) (*user_domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}

type fakeOrganizationRepo struct {
	organizations map[string]*domain.Organization
}

func (r *fakeOrganizationRepo) FindById(id string) (*domain.Organization, error) {
	org, ok := r.organizations[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return org, nil
}

type fakeTeamRepo struct {
	teams map[string]*domain.Team
}

func (r *fakeTeamRepo) FindById(id string) (*domain.Team, error) {
	team, ok := r.teams[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return team, nil
}

type fakeLeagueRepo struct {
	mu     sync.Mutex
	league *domain.League
}

func (r *fakeLeagueRepo) FindById(id string) (*domain.League, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return cloneLeague(r.league), nil
}

func (r *fakeLeagueRepo) Save(*domain.League) error {
	return nil
}

func (r *fakeLeagueRepo) Search(dtos.LeagueSearchDTO) (*map[string]interface{}, error) {
	return nil, nil
}

func (r *fakeLeagueRepo) FetchLeagueMembers(string) ([]interface{}, error) {
	return nil, nil
}

func (r *fakeLeagueRepo) FetchLeagueDetails(domain.League) (map[string]interface{}, error) {
	return nil, nil
}

func (r *fakeLeagueRepo) WithLockedLeague(_ string, fn func(*domain.League) error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	working := cloneLeague(r.league)
	if err := fn(working); err != nil {
		return err
	}

	r.league = working
	return nil
}

type fakeSeasonRepo struct {
	mu                sync.Mutex
	season            *domain.Season
	findBarrierTarget int
	findBarrierCount  int
	findBarrier       chan struct{}
}

func (r *fakeSeasonRepo) FindByID(string) (*domain.Season, error) {
	r.mu.Lock()
	season := cloneSeason(r.season)
	if r.findBarrierTarget > 0 {
		r.findBarrierCount++
		if r.findBarrierCount == r.findBarrierTarget {
			close(r.findBarrier)
		}
		barrier := r.findBarrier
		r.mu.Unlock()
		<-barrier
		return season, nil
	}
	r.mu.Unlock()
	return season, nil
}

func (r *fakeSeasonRepo) Save(season *domain.Season) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if season.Version <= 0 {
		r.season = cloneSeason(season)
		r.season.Version = 1
		return nil
	}

	if r.season.Version != season.Version {
		return app_errors.ErrConcurrentModification
	}

	next := cloneSeason(season)
	next.Version++
	r.season = next
	return nil
}

func (r *fakeSeasonRepo) Search(string, string, dtos.SearchSeasonDTO) ([]interface{}, int) {
	return nil, 0
}

func (r *fakeSeasonRepo) FetchDetails(string) (map[string]interface{}, error) {
	return nil, nil
}

func (r *fakeSeasonRepo) FetchSeasonStandings(string) (map[string]interface{}, error) {
	return nil, nil
}

func (r *fakeSeasonRepo) FetchSeasonMatchUps(string) ([]interface{}, error) {
	return nil, nil
}

type fakeSeasonLeagueRepo struct {
	league *domain.League
}

func (r *fakeSeasonLeagueRepo) FindById(string) (*domain.League, error) {
	return cloneLeague(r.league), nil
}

func TestLeagueService_InitiateTeamMembership_ConcurrentAddsKeepBothMemberships(t *testing.T) {
	ownerID := "owner-1"
	orgID := "org-1"
	leagueID := "league-1"
	team1ID := domain.TeamId("team-1")
	team2ID := domain.TeamId("team-2")

	service := &LeagueService{
		leagueRepository: &fakeLeagueRepo{
			league: &domain.League{
				Id:             stringPtr(leagueID),
				Name:           "Summer League",
				OwnerId:        ownerID,
				OrganizationId: orgID,
				Memberships:    []domain.LeagueMembership{},
			},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*domain.Organization{
				orgID: {ID: stringPtr(orgID), Name: "Org", OrganizationOwnerId: ownerID},
			},
		},
		teamRepository: &fakeTeamRepo{
			teams: map[string]*domain.Team{
				string(team1ID): {ID: &team1ID, Name: "Team 1", OrganizationId: orgID},
				string(team2ID): {ID: &team2ID, Name: "Team 2", OrganizationId: orgID},
			},
		},
		userRepo: &fakeUserRepo{
			users: map[string]*user_domain.User{
				ownerID: {Id: ownerID, Email: "owner@example.com"},
			},
		},
	}

	var wg sync.WaitGroup
	errs := make(chan error, 2)

	for _, teamID := range []string{string(team1ID), string(team2ID)} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			errs <- service.InitiateTeamMembership(ownerID, leagueID, id)
		}(teamID)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("expected both concurrent membership commands to succeed, got %v", err)
		}
	}

	repo := service.leagueRepository.(*fakeLeagueRepo)
	if len(repo.league.Memberships) != 2 {
		t.Fatalf("expected 2 memberships after concurrent adds, got %d", len(repo.league.Memberships))
	}
}

func TestSeasonService_StartSeason_ConcurrentStartConflictsOnVersion(t *testing.T) {
	ownerID := "owner-1"
	orgID := "org-1"
	leagueID := "league-1"
	seasonID := "season-1"

	seasonRepo := &fakeSeasonRepo{
		season: &domain.Season{
			ID:       seasonID,
			LeagueId: leagueID,
			Name:     "Spring",
			Status:   domain.SeasonStatusPlanned,
			Version:  1,
			Rounds: []domain.Round{
				{RoundNumber: 1, Matches: []domain.Match{{ID: "match-1", HomeTeamID: "team-1", AwayTeamID: "team-2"}}},
			},
		},
		findBarrierTarget: 2,
		findBarrier:       make(chan struct{}),
	}

	service := &SeasonService{
		seasonRepository: seasonRepo,
		leagueRepository: &fakeSeasonLeagueRepo{
			league: &domain.League{Id: stringPtr(leagueID), OrganizationId: orgID},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*domain.Organization{
				orgID: {ID: stringPtr(orgID), Name: "Org", OrganizationOwnerId: ownerID},
			},
		},
	}

	var wg sync.WaitGroup
	errs := make(chan error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- service.StartSeason(ownerID, seasonID)
		}()
	}

	wg.Wait()
	close(errs)

	successes := 0
	conflicts := 0
	for err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, app_errors.ErrConcurrentModification):
			conflicts++
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if successes != 1 || conflicts != 1 {
		t.Fatalf("expected one success and one concurrency conflict, got successes=%d conflicts=%d", successes, conflicts)
	}
	if seasonRepo.season.Status != domain.SeasonStatusInProgress {
		t.Fatalf("expected final season status to be in_progress, got %s", seasonRepo.season.Status)
	}
}

func TestSeasonService_ChangeMatchScore_ConcurrentUpdatesConflictOnVersion(t *testing.T) {
	ownerID := "owner-1"
	orgID := "org-1"
	leagueID := "league-1"
	seasonID := "season-1"
	matchID := "match-1"

	seasonRepo := &fakeSeasonRepo{
		season: &domain.Season{
			ID:       seasonID,
			LeagueId: leagueID,
			Name:     "Spring",
			Status:   domain.SeasonStatusInProgress,
			Version:  1,
			Rounds: []domain.Round{
				{RoundNumber: 1, Matches: []domain.Match{{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2"}}},
			},
		},
		findBarrierTarget: 2,
		findBarrier:       make(chan struct{}),
	}

	service := &SeasonService{
		seasonRepository: seasonRepo,
		leagueRepository: &fakeSeasonLeagueRepo{
			league: &domain.League{Id: stringPtr(leagueID), OrganizationId: orgID},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*domain.Organization{
				orgID: {ID: stringPtr(orgID), Name: "Org", OrganizationOwnerId: ownerID},
			},
		},
	}

	var wg sync.WaitGroup
	errs := make(chan error, 2)

	for _, score := range []int{1, 2} {
		wg.Add(1)
		go func(homeScore int) {
			defer wg.Done()
			errs <- service.ChangeMatchUpScore(ownerID, seasonID, dtos.ChangeGameScoreDTO{
				MatchID:   matchID,
				HomeScore: homeScore,
				AwayScore: 0,
			})
		}(score)
	}

	wg.Wait()
	close(errs)

	successes := 0
	conflicts := 0
	for err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, app_errors.ErrConcurrentModification):
			conflicts++
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if successes != 1 || conflicts != 1 {
		t.Fatalf("expected one success and one concurrency conflict, got successes=%d conflicts=%d", successes, conflicts)
	}
}

func cloneLeague(in *domain.League) *domain.League {
	if in == nil {
		return nil
	}

	memberships := make([]domain.LeagueMembership, len(in.Memberships))
	copy(memberships, in.Memberships)

	out := *in
	out.Memberships = memberships
	return &out
}

func cloneSeason(in *domain.Season) *domain.Season {
	if in == nil {
		return nil
	}

	rounds := make([]domain.Round, len(in.Rounds))
	for i, round := range in.Rounds {
		matches := make([]domain.Match, len(round.Matches))
		copy(matches, round.Matches)
		rounds[i] = domain.Round{
			RoundNumber: round.RoundNumber,
			Matches:     matches,
		}
	}

	out := *in
	out.Rounds = rounds
	return &out
}

func stringPtr(value string) *string {
	return &value
}

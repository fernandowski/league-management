package services

import (
	"errors"
	leaguepkg "league-management/internal/organization_management/domain/league"
	organizationpkg "league-management/internal/organization_management/domain/organization"
	seasonpkg "league-management/internal/organization_management/domain/season"
	teampkg "league-management/internal/organization_management/domain/team"
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
	organizations map[string]*organizationpkg.Organization
}

func (r *fakeOrganizationRepo) FindById(id string) (*organizationpkg.Organization, error) {
	org, ok := r.organizations[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return org, nil
}

type fakeTeamRepo struct {
	teams map[string]*teampkg.Team
}

func (r *fakeTeamRepo) FindById(id string) (*teampkg.Team, error) {
	team, ok := r.teams[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return team, nil
}

type fakeLeagueRepo struct {
	mu     sync.Mutex
	league *leaguepkg.League
}

func (r *fakeLeagueRepo) FindById(id string) (*leaguepkg.League, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return cloneLeague(r.league), nil
}

func (r *fakeLeagueRepo) Save(*leaguepkg.League) error {
	return nil
}

func (r *fakeLeagueRepo) Search(dtos.LeagueSearchDTO) (*map[string]interface{}, error) {
	return nil, nil
}

func (r *fakeLeagueRepo) FetchLeagueMembers(string) ([]interface{}, error) {
	return nil, nil
}

func (r *fakeLeagueRepo) FetchLeagueDetails(leaguepkg.League) (map[string]interface{}, error) {
	return nil, nil
}

func (r *fakeLeagueRepo) WithLockedLeague(_ string, fn func(*leaguepkg.League) error) error {
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
	season            *seasonpkg.Season
	findBarrierTarget int
	findBarrierCount  int
	findBarrier       chan struct{}
}

func (r *fakeSeasonRepo) FindByID(string) (*seasonpkg.Season, error) {
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

func (r *fakeSeasonRepo) Save(season *seasonpkg.Season) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if season.CurrentVersion() <= 0 {
		r.season = cloneSeason(season)
		snapshot := r.season.Snapshot()
		snapshot.Version = 1
		r.season = seasonpkg.RehydrateSeasonFromSnapshot(snapshot)
		r.season.ApplyPersistedVersion(1)
		return nil
	}

	if r.season.CurrentVersion() != season.CurrentVersion() {
		return app_errors.ErrConcurrentModification
	}

	next := cloneSeason(season)
	snapshot := next.Snapshot()
	snapshot.Version = next.CurrentVersion() + 1
	r.season = seasonpkg.RehydrateSeasonFromSnapshot(snapshot)
	r.season.ApplyPersistedVersion(next.CurrentVersion() + 1)
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

func (r *fakeSeasonRepo) FetchPlayoffBracket(string) (map[string]interface{}, error) {
	return nil, nil
}

type fakeSeasonLeagueRepo struct {
	league *leaguepkg.League
}

func (r *fakeSeasonLeagueRepo) FindById(string) (*leaguepkg.League, error) {
	return cloneLeague(r.league), nil
}

func TestLeagueService_InitiateTeamMembership_ConcurrentAddsKeepBothMemberships(t *testing.T) {
	ownerID := "owner-1"
	orgID := "org-1"
	leagueID := "league-1"
	team1ID := teampkg.TeamId("team-1")
	team2ID := teampkg.TeamId("team-2")

	service := &LeagueService{
		leagueRepository: &fakeLeagueRepo{
			league: &leaguepkg.League{
				Id:             stringPtr(leagueID),
				Name:           "Summer League",
				OwnerId:        ownerID,
				OrganizationId: orgID,
				Memberships:    []leaguepkg.Membership{},
			},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*organizationpkg.Organization{
				orgID: {ID: stringPtr(orgID), Name: "Org", OrganizationOwnerId: ownerID},
			},
		},
		teamRepository: &fakeTeamRepo{
			teams: map[string]*teampkg.Team{
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
		season: rehydratedSeason(
			seasonID,
			leagueID,
			"Spring",
			seasonpkg.SeasonStatusPlanned,
			"",
			1,
			[]seasonpkg.Round{
				seasonpkg.RehydrateRound(1, []seasonpkg.Match{seasonpkg.RehydrateMatch(seasonpkg.MatchState{ID: "match-1", HomeTeamID: "team-1", AwayTeamID: "team-2"})}),
			},
			nil,
			nil,
			nil,
		),
		findBarrierTarget: 2,
		findBarrier:       make(chan struct{}),
	}

	service := &SeasonService{
		seasonRepository: seasonRepo,
		leagueRepository: &fakeSeasonLeagueRepo{
			league: &leaguepkg.League{Id: stringPtr(leagueID), OrganizationId: orgID},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*organizationpkg.Organization{
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
	if seasonRepo.season.CurrentStatus() != seasonpkg.SeasonStatusInProgress {
		t.Fatalf("expected final season status to be in_progress, got %s", seasonRepo.season.CurrentStatus())
	}
}

func TestSeasonService_ChangeMatchScore_ConcurrentUpdatesConflictOnVersion(t *testing.T) {
	ownerID := "owner-1"
	orgID := "org-1"
	leagueID := "league-1"
	seasonID := "season-1"
	matchID := "match-1"

	seasonRepo := &fakeSeasonRepo{
		season: rehydratedSeason(
			seasonID,
			leagueID,
			"Spring",
			seasonpkg.SeasonStatusInProgress,
			"",
			1,
			[]seasonpkg.Round{
				seasonpkg.RehydrateRound(1, []seasonpkg.Match{seasonpkg.RehydrateMatch(seasonpkg.MatchState{ID: matchID, HomeTeamID: "team-1", AwayTeamID: "team-2", Status: seasonpkg.MatchStatusInProgress})}),
			},
			nil,
			nil,
			nil,
		),
		findBarrierTarget: 2,
		findBarrier:       make(chan struct{}),
	}

	service := &SeasonService{
		seasonRepository: seasonRepo,
		leagueRepository: &fakeSeasonLeagueRepo{
			league: &leaguepkg.League{Id: stringPtr(leagueID), OrganizationId: orgID},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*organizationpkg.Organization{
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

func TestSeasonService_ConfigurePlayoffRules_FinishedRegularSeasonAllowsOwnerToSave(t *testing.T) {
	ownerID := "owner-1"
	orgID := "org-1"
	leagueID := "league-1"
	seasonID := "season-1"

	seasonRepo := &fakeSeasonRepo{
		season: rehydratedSeason(
			seasonID,
			leagueID,
			"Spring",
			seasonpkg.SeasonStatusFinished,
			seasonpkg.SeasonPhaseRegularSeason,
			1,
			nil,
			nil,
			nil,
			nil,
		),
	}

	service := &SeasonService{
		seasonRepository: seasonRepo,
		leagueRepository: &fakeSeasonLeagueRepo{
			league: &leaguepkg.League{Id: stringPtr(leagueID), OrganizationId: orgID},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*organizationpkg.Organization{
				orgID: {ID: stringPtr(orgID), Name: "Org", OrganizationOwnerId: ownerID},
			},
		},
	}

	err := service.ConfigurePlayoffRules(ownerID, seasonID, dtos.ConfigurePlayoffRulesDTO{
		QualificationType: "top_n",
		QualifierCount:    4,
		Rounds: []dtos.ConfigurePlayoffRoundDTO{
			{Name: "semifinal", Legs: 2, TiedAggregateResolution: "penalties"},
			{Name: "final", Legs: 1, TiedAggregateResolution: "penalties"},
		},
	})

	if err != nil {
		t.Fatalf("expected playoff rules to save for finished regular season, got %v", err)
	}
	qualifierCount, configured := seasonRepo.season.PlayoffQualifierCount()
	if !configured || !seasonRepo.season.HasPlayoffRules() {
		t.Fatalf("expected playoff rules to be stored on season")
	}
	if qualifierCount != 4 {
		t.Fatalf("expected qualifier count 4, got %d", qualifierCount)
	}
}

func TestSeasonService_ConfigurePlayoffRules_BeforePlayoffMatchPlayedClearsBracket(t *testing.T) {
	ownerID := "owner-1"
	orgID := "org-1"
	leagueID := "league-1"
	seasonID := "season-1"

	seasonRepo := &fakeSeasonRepo{
		season: rehydratedSeason(
			seasonID,
			leagueID,
			"Spring",
			seasonpkg.SeasonStatusInProgress,
			seasonpkg.SeasonPhasePlayoffs,
			1,
			nil,
			nil,
			seasonpkg.RehydratePlayoffBracket([]seasonpkg.PlayoffBracketRoundSnapshot{
				{
					Name:  "semifinal",
					Order: 1,
					Ties: []seasonpkg.PlayoffTieSnapshot{
						{
							ID:     "tie-1",
							Status: "ready",
							Matches: []seasonpkg.MatchSnapshot{
								seasonpkg.RehydrateMatch(seasonpkg.MatchState{ID: "match-1", PlayoffTieID: "tie-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: seasonpkg.MatchStatusScheduled}).Snapshot(),
							},
						},
					},
				},
			}),
			nil,
		),
	}

	service := &SeasonService{
		seasonRepository: seasonRepo,
		leagueRepository: &fakeSeasonLeagueRepo{
			league: &leaguepkg.League{Id: stringPtr(leagueID), OrganizationId: orgID},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*organizationpkg.Organization{
				orgID: {ID: stringPtr(orgID), Name: "Org", OrganizationOwnerId: ownerID},
			},
		},
	}

	err := service.ConfigurePlayoffRules(ownerID, seasonID, dtos.ConfigurePlayoffRulesDTO{
		QualificationType: "top_n",
		QualifierCount:    2,
		Rounds: []dtos.ConfigurePlayoffRoundDTO{
			{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
		},
	})

	if err != nil {
		t.Fatalf("expected playoff rules to save before playoff matches are played, got %v", err)
	}
	if seasonRepo.season.HasPlayoffBracket() {
		t.Fatalf("expected existing playoff bracket to be cleared after rules change")
	}
	if seasonRepo.season.CurrentPhase() != seasonpkg.SeasonPhaseRegularSeason {
		t.Fatalf("expected season phase to return to regular_season, got %s", seasonRepo.season.CurrentPhase())
	}
}

func TestSeasonService_ConfigurePlayoffRules_RejectsNonOwner(t *testing.T) {
	ownerID := "owner-1"
	otherUserID := "user-2"
	orgID := "org-1"
	leagueID := "league-1"
	seasonID := "season-1"

	seasonRepo := &fakeSeasonRepo{
		season: rehydratedSeason(
			seasonID,
			leagueID,
			"Spring",
			seasonpkg.SeasonStatusFinished,
			seasonpkg.SeasonPhaseRegularSeason,
			1,
			nil,
			nil,
			nil,
			nil,
		),
	}

	service := &SeasonService{
		seasonRepository: seasonRepo,
		leagueRepository: &fakeSeasonLeagueRepo{
			league: &leaguepkg.League{Id: stringPtr(leagueID), OrganizationId: orgID},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*organizationpkg.Organization{
				orgID: {ID: stringPtr(orgID), Name: "Org", OrganizationOwnerId: ownerID},
			},
		},
	}

	err := service.ConfigurePlayoffRules(otherUserID, seasonID, dtos.ConfigurePlayoffRulesDTO{
		QualificationType: "top_n",
		QualifierCount:    4,
		Rounds: []dtos.ConfigurePlayoffRoundDTO{
			{Name: "semifinal", Legs: 2, TiedAggregateResolution: "penalties"},
		},
	})

	if err == nil {
		t.Fatalf("expected non-owner playoff configuration to fail")
	}
}

func TestSeasonService_ConfigurePlayoffRules_RejectsTwoLegFinal(t *testing.T) {
	ownerID := "owner-1"
	orgID := "org-1"
	leagueID := "league-1"
	seasonID := "season-1"

	seasonRepo := &fakeSeasonRepo{
		season: rehydratedSeason(
			seasonID,
			leagueID,
			"Spring",
			seasonpkg.SeasonStatusFinished,
			seasonpkg.SeasonPhaseRegularSeason,
			1,
			nil,
			nil,
			nil,
			nil,
		),
	}

	service := &SeasonService{
		seasonRepository: seasonRepo,
		leagueRepository: &fakeSeasonLeagueRepo{
			league: &leaguepkg.League{Id: stringPtr(leagueID), OrganizationId: orgID},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*organizationpkg.Organization{
				orgID: {ID: stringPtr(orgID), Name: "Org", OrganizationOwnerId: ownerID},
			},
		},
	}

	err := service.ConfigurePlayoffRules(ownerID, seasonID, dtos.ConfigurePlayoffRulesDTO{
		QualificationType: "top_n",
		QualifierCount:    4,
		Rounds: []dtos.ConfigurePlayoffRoundDTO{
			{Name: "semifinal", Legs: 2, TiedAggregateResolution: "higher_seed_advances"},
			{Name: "final", Legs: 2, TiedAggregateResolution: "clear_winner_required"},
		},
	})

	if err == nil {
		t.Fatalf("expected two-leg final playoff configuration to fail")
	}
}

func TestSeasonService_PlanSchedule_WithActiveLeagueMembershipsSucceeds(t *testing.T) {
	ownerID := "owner-1"
	orgID := "org-1"
	leagueID := "league-1"
	seasonID := "season-1"

	seasonRepo := &fakeSeasonRepo{
		season: rehydratedSeason(
			seasonID,
			leagueID,
			"Spring",
			seasonpkg.SeasonStatusPending,
			seasonpkg.SeasonPhaseRegularSeason,
			1,
			nil,
			nil,
			nil,
			nil,
		),
	}

	service := &SeasonService{
		seasonRepository: seasonRepo,
		leagueRepository: &fakeSeasonLeagueRepo{
			league: &leaguepkg.League{
				Id:             stringPtr(leagueID),
				OrganizationId: orgID,
				Memberships: []leaguepkg.Membership{
					{ID: "m1", TeamID: "team-1", MemberShipStatus: leaguepkg.MembershipActive},
					{ID: "m2", TeamID: "team-2", MemberShipStatus: leaguepkg.MembershipActive},
					{ID: "m3", TeamID: "team-3", MemberShipStatus: leaguepkg.MembershipActive},
				},
			},
		},
		organizationRepo: &fakeOrganizationRepo{
			organizations: map[string]*organizationpkg.Organization{
				orgID: {ID: stringPtr(orgID), Name: "Org", OrganizationOwnerId: ownerID},
			},
		},
	}

	err := service.PlanSchedule(ownerID, seasonID)
	if err != nil {
		t.Fatalf("expected planning to succeed with active memberships, got %v", err)
	}

	if len(seasonRepo.season.PlannedRounds()) == 0 {
		t.Fatalf("expected planned rounds to be created")
	}
	if seasonRepo.season.CurrentStatus() != seasonpkg.SeasonStatusPlanned {
		t.Fatalf("expected season status to be planned, got %s", seasonRepo.season.CurrentStatus())
	}
}

func cloneLeague(in *leaguepkg.League) *leaguepkg.League {
	if in == nil {
		return nil
	}

	memberships := make([]leaguepkg.Membership, len(in.Memberships))
	copy(memberships, in.Memberships)

	out := *in
	out.Memberships = memberships
	return &out
}

func cloneSeason(in *seasonpkg.Season) *seasonpkg.Season {
	if in == nil {
		return nil
	}

	return seasonpkg.RehydrateSeasonFromSnapshot(in.Snapshot())
}

func rehydratedSeason(
	id, leagueID, name string,
	status seasonpkg.SeasonStatus,
	phase seasonpkg.SeasonPhase,
	version int,
	rounds []seasonpkg.Round,
	playoffRules *seasonpkg.PlayoffRulesSnapshot,
	playoffBracket *seasonpkg.PlayoffBracketSnapshot,
	championTeamID *string,
) *seasonpkg.Season {
	roundSnapshots := make([]seasonpkg.RoundSnapshot, len(rounds))
	for i, round := range rounds {
		roundSnapshots[i] = round.Snapshot()
	}

	return seasonpkg.RehydrateSeasonFromSnapshot(seasonpkg.SeasonSnapshot{
		ID:             id,
		LeagueID:       leagueID,
		Name:           name,
		Status:         status,
		Phase:          phase,
		Version:        version,
		Rounds:         roundSnapshots,
		PlayoffRules:   playoffRules,
		PlayoffBracket: playoffBracket,
		ChampionTeamID: championTeamID,
	})
}

func stringPtr(value string) *string {
	return &value
}

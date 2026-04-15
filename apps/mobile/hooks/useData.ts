import {apiRequest} from "@/api/api";
import {useCallback, useState} from "react";


export interface LeagueMembershipResponse {
    league_id: string,
    membership_id: string,
    team_id: string,
    team_name: string,
}

export interface TeamResponse {
    id: string
    name: string
}

export interface LeagueDetailResponse {
    id: string,
    name: string
    season: Record<string, any>
    active_members: number
}


export interface Match {
    home_team: string,
    away_team: string,
    id: string
}

export interface Round {
    round_number: string
    matches: Match[]
}
export interface SeasonDetailResponse {
    id: string,
    name: string
    status: string
    phase?: string
    champion_team_id?: string | null
    rounds: Round[]
}

export interface PlayoffRoundRuleResponse {
    name: string
    legs: number
    higher_seed_hosts_second_leg: boolean
    tied_aggregate_resolution: string
}

export interface PlayoffRulesResponse {
    season_id: string
    season_status: string
    season_phase: string
    configured: boolean
    bracket_generated: boolean
    playoffs_started: boolean
    rules_locked: boolean
    rules: {
        qualification_type: string
        qualifier_count: number
        reseed_each_round: boolean
        third_place_match: boolean
        allow_admin_seed_override: boolean
        rounds: PlayoffRoundRuleResponse[]
    } | null
}

export interface PlayoffQualifiedTeamResponse {
    team_id: string
    team_name: string
    seed: number
    rank: number
    points: number
    wins: number
    losses: number
    ties: number
    total_goals: number
}

export interface PlayoffQualificationPreviewResponse {
    season_id: string
    qualifier_count: number
    bracket_exists: boolean
    qualified_teams: PlayoffQualifiedTeamResponse[]
}

export interface PlayoffBracketTeamResponse {
    id: string | null
    name: string
}

export interface PlayoffTieResponse {
    id: string
    slot_order: number
    status: string
    home_seed: number | null
    away_seed: number | null
    winner_team_id: string | null
    home_team: PlayoffBracketTeamResponse
    away_team: PlayoffBracketTeamResponse
    matches: PlayoffTieMatchResponse[]
}

export interface PlayoffTieMatchResponse {
    id: string
    match_order: number
    home_score: number
    away_score: number
    status: string
    home_team: string
    away_team: string
}

export interface PlayoffBracketRoundResponse {
    name: string
    order: number
    ties: PlayoffTieResponse[]
}

export interface PlayoffBracketResponse {
    season_id: string
    generated: boolean
    rounds: PlayoffBracketRoundResponse[]
}

export interface SeasonStandings {
    games_played: number
    team_id: string
    team_name: string
    total_goals: number
    total_losses: number
    total_points: number
    total_ties: number
    total_wins: number
}
export interface SeasonStandingsResponse {
    season_id: string,
    standings: SeasonStandings[]
}


export interface MatchScore {
    id: string
    away_score: number
    home_score: number
    away_team: string
    home_team: string
    round: number
    status: string
}
export interface MatchesResponse {
    round: number
    matches: MatchScore[]
}

export function useData<TResponse>() {
    const [fetching, setFetching] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [data, setData] = useState<TResponse | null>();
    const fetchData = useCallback(async (endpoint: string): Promise<TResponse | undefined> => {
        setError(null);
        setFetching(true);
        try {
            const response: TResponse = await apiRequest(endpoint, {method: "GET"});
            setFetching(false)
            setData(response);
            return response;
        } catch (e: any) {
            setError(e.message || "Error Fetching data.");
            setFetching(false);
        }
    }, []);

    return { fetchData, fetching, error, data };
}

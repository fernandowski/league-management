import {apiRequest} from "@/api/api";
import {useState} from "react";


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
    rounds: Round[]
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
    const fetchData = async (endpoint: string): Promise<TResponse | undefined> => {
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
    }

    return { fetchData, fetching, error, data };
}

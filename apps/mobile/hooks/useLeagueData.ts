import {useState} from "react";
import {apiRequest} from "@/api/api";

export interface LeagueSearchParams {
    organization_id: string | null
    limit: number
    offset: number
    term: string | ''
}

export interface League {
    id: string
    name: string
    teamIds: string[],
    totalMembers: number
}


export interface LeagueAPIResponse {
    id: string
    name: string
    total_members: number
}
export interface LeagueResponse {
    data: LeagueAPIResponse[]
    total: number
}

export function useLeagueData(): {
    fetchData: (params: LeagueSearchParams) => Promise<void>,
    fetching: boolean,
    error: null | string,
    data: League[],
    total: number
} {
    const [fetching, setFetching] = useState(false);
    const [error, setError] = useState<null | string>(null);
    const [data, setData] = useState<League[]>([]);
    const [total, setTotal] = useState(0);

    const fetchData = async (params: LeagueSearchParams): Promise<void> => {
        setFetching(true);
        try {
            const offsetQuery = `offset=${params.offset}`;
            const limitQuery = `limit=${params.limit}`;
            const termQuery = params.term.trim() !== '' ? `term=${params.term}` : '';
            const organizationQuery = `organization_id=${params.organization_id}`;
            const response: LeagueResponse = await apiRequest(
                `/v1/leagues?${[offsetQuery, limitQuery, termQuery, organizationQuery].join('&')}`,
                {method: "GET"}
            );
            setFetching(false);
            const leagues: League[] = response.data.map((league: LeagueAPIResponse) => ({
                id: league.id,
                name: league.name,
                teamIds: [],
                totalMembers: league.total_members
            }))

            setData(leagues)
            setTotal(response.total)
        } catch (e) {
            setError("Error Fetching leagues");
            setFetching(false);
        }
    };

    return {fetchData, fetching, error, data, total};
}

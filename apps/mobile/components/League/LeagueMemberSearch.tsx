import {useEffect, useState} from "react";
import {Button, DataTable, Searchbar, Text} from "react-native-paper";
import {View} from "react-native";
import {TeamResponse, useData} from "@/hooks/useData";
import {useDebounce} from "@/hooks/useDebounce";
import {apiRequest} from "@/api/api";
import {LeagueMember} from "@/components/League/LeagueMembers";

export interface LeagueMemberSearchProps {
    organizationId: string
    leagueId: string
    onTeamAdded: () => void
    members: LeagueMember[]
}

export interface Team {
    id: string
    name: string
}

export default function LeagueMemberSearch(props: LeagueMemberSearchProps) {
    const [searchQuery, setSearchQuery] = useState('');
    const {fetchData, fetching, error} = useData<TeamResponse[]>();
    const [teams, setTeams] = useState<Team[]>([]);
    const debouncedSearchQuery = useDebounce(searchQuery, 800);

    const fetchTeams = async (term: string) => {
        const encodedTerm: string = encodeURIComponent(`%${term}%`)
        const response: TeamResponse[] | undefined = await fetchData(`/v1/teams/?organization_id=${props.organizationId}&term=${encodedTerm}`);
        if (response) {
            const teams: Team[] = response.map((team: Team) => ({id: team.id, name: team.name}));
            setTeams(teams)
        }
    }

    const handleAdd = async (teamId: string) => {
        try {
            await apiRequest(
                `/v1/leagues/${props.leagueId}/invites`,
                {
                    method: 'POST',
                    body: {team_id: teamId}
                },
            );
            props.onTeamAdded();
        } catch {

        }
    }

    useEffect(() => {
        if (debouncedSearchQuery !== '' && !fetching) {
            fetchTeams(debouncedSearchQuery)
        }
    }, [debouncedSearchQuery]);

    return (
        <View>
            <Searchbar
                placeholder="Search"
                onChangeText={setSearchQuery}
                value={searchQuery}
            />
            <DataTable>
                {
                    teams.map((team: Team) => {
                        return (
                            <DataTable.Row key={team.id}>
                                <DataTable.Cell style={{flex: 5}}>{team.name}</DataTable.Cell>
                                <DataTable.Cell style={{flex: 1}}>
                                    <Button
                                        mode={"elevated"}
                                        onPress={() => handleAdd(team.id)}
                                        disabled={isMember(team.id, props.members)}
                                    >
                                        Add
                                    </Button>
                                </DataTable.Cell>
                            </DataTable.Row>
                        )
                    })
                }
            </DataTable>
        </View>
    )
}

function isMember (memberId: string, members: LeagueMember[]): boolean {
    const index = members.findIndex(member => (member.teamId === memberId));
    return index >= 0;
}

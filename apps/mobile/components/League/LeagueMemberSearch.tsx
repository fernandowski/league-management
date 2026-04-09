import {useCallback, useEffect, useState} from "react";
import {Button, Searchbar} from "react-native-paper";
import {StyleSheet, View} from "react-native";
import {TeamResponse, useData} from "@/hooks/useData";
import {useDebounce} from "@/hooks/useDebounce";
import {apiRequest} from "@/api/api";
import {LeagueMember} from "@/components/League/LeagueMembers";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";

export interface LeagueMemberSearchProps {
    organizationId: string
    leagueId: string
    onTeamAdded: () => void
    members: LeagueMember[]
    onCreateTeam: (teamName: string) => void
    creatingTeam?: boolean
    createError?: string | null
}

export interface Team {
    id: string
    name: string
}

export default function LeagueMemberSearch(props: LeagueMemberSearchProps) {
    const [searchQuery, setSearchQuery] = useState('');
    const {fetchData, fetching} = useData<TeamResponse[]>();
    const [teams, setTeams] = useState<Team[]>([]);
    const debouncedSearchQuery = useDebounce(searchQuery, 800);
    const theme = useAppTheme();

    const fetchTeams = useCallback(async (term: string) => {
        const encodedTerm: string = encodeURIComponent(`%${term}%`)
        const response: TeamResponse[] | undefined = await fetchData(`/v1/teams/?organization_id=${props.organizationId}&term=${encodedTerm}`);
        if (response) {
            const teams: Team[] = response.map((team: Team) => ({id: team.id, name: team.name}));
            setTeams(teams)
        }
    }, [fetchData, props.organizationId])

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
        if (debouncedSearchQuery !== '') {
            fetchTeams(debouncedSearchQuery)
        } else if (debouncedSearchQuery === '') {
            setTeams([]);
        }
    }, [debouncedSearchQuery, fetchTeams]);

    return (
        <View style={styles.container}>
            <Searchbar
                placeholder="Search teams by name"
                onChangeText={setSearchQuery}
                value={searchQuery}
            />
            {searchQuery.trim() === '' ? (
                <View style={[styles.helperCard, {backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant}]}>
                    <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                        Start typing to find teams available to add.
                    </AppText>
                </View>
            ) : teams.length === 0 && !fetching ? (
                <View style={[styles.helperCard, {backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant}]}>
                    <View style={styles.emptyStateCopy}>
                        <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                            No teams matched this search.
                        </AppText>
                        <AppText variant="bodySmall" style={{color: theme.colors.onSurfaceVariant}}>
                            Create <AppText style={{fontWeight: "600", color: theme.colors.onSurface}}>{searchQuery.trim()}</AppText> and add it to this league.
                        </AppText>
                        {props.createError ? (
                            <AppText variant="bodySmall" style={{color: theme.colors.error}}>
                                {props.createError}
                            </AppText>
                        ) : null}
                    </View>
                    <Button
                        mode="contained"
                        onPress={() => props.onCreateTeam(searchQuery)}
                        loading={props.creatingTeam}
                        disabled={props.creatingTeam}
                    >
                        Create team
                    </Button>
                </View>
            ) : (
                <View style={styles.resultsList}>
                    {teams.map((team: Team) => {
                        const memberAlreadyAdded = isMember(team.id, props.members);

                        return (
                            <View
                                key={team.id}
                                style={[
                                    styles.resultRow,
                                    {
                                        backgroundColor: theme.colors.surface,
                                        borderColor: theme.colors.outlineVariant,
                                    },
                                ]}
                            >
                                <View style={styles.resultMeta}>
                                    <AppText variant="titleMedium">{team.name}</AppText>
                                    <AppText variant="bodySmall" style={{color: theme.colors.onSurfaceVariant}}>
                                        {memberAlreadyAdded ? "Already in this league" : "Available to add"}
                                    </AppText>
                                </View>
                                <Button
                                    mode={"contained"}
                                    onPress={() => handleAdd(team.id)}
                                    disabled={memberAlreadyAdded}
                                >
                                    {memberAlreadyAdded ? "Added" : "Add"}
                                </Button>
                            </View>
                        )
                    })}
                </View>
            )}
        </View>
    )
}

function isMember (memberId: string, members: LeagueMember[]): boolean {
    const index = members.findIndex(member => (member.teamId === memberId));
    return index >= 0;
}

const styles = StyleSheet.create({
    container: {
        gap: 12,
    },
    helperCard: {
        borderWidth: 1,
        borderRadius: 18,
        padding: 16,
        gap: 12,
    },
    emptyStateCopy: {
        gap: 4,
    },
    resultsList: {
        gap: 10,
    },
    resultRow: {
        borderWidth: 1,
        borderRadius: 18,
        padding: 14,
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
        gap: 12,
    },
    resultMeta: {
        flex: 1,
        gap: 2,
    },
});

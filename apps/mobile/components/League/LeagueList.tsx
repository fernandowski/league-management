import {Animated, View, StyleSheet} from "react-native";
import ScrollView = Animated.ScrollView;
import {League, useLeagueData} from "@/hooks/useLeagueData";
import React, {useCallback, useEffect, useState} from "react";
import {useOrganizationStore} from "@/stores/organizationStore";
import {Link} from "expo-router";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import Pagination from "@/components/Pagination/Pagination";
import {useAppTheme} from "@/theme/theme";

interface LeagueListProps {
    refresh?: boolean
}

export function LeagueList(props: LeagueListProps): React.JSX.Element {
    const {fetchData, data, total} = useLeagueData();
    const {organization} = useOrganizationStore()
    const theme = useAppTheme();
    const [page, setPage] = useState(0);
    const itemsPerPage = 6;

    const fetchLeagues = useCallback(async (nextPage: number) => {
        if (organization !== null) {
            await fetchData({
                organization_id: organization,
                limit: itemsPerPage,
                offset: nextPage * itemsPerPage,
                term: "",
            });
        }
    }, [fetchData, itemsPerPage, organization]);

    useEffect(() => {
        if (organization !== null) {
            setPage(0);
            fetchLeagues(0);
        }

    }, [fetchLeagues, organization, props.refresh]);

    useEffect(() => {
        if (organization !== null) {
            fetchLeagues(page);
        }
    }, [fetchLeagues, organization, page]);

    return (
        <View style={{flex: 1}}>
            {data.length === 0 ? (
                <AppCard style={[styles.emptyCard, { backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant }]}>
                    <AppCard.Content style={styles.emptyContent}>
                        <AppText style={[styles.title, {color: theme.colors.onSurface}]}>You do not have any leagues</AppText>
                        <AppText style={[styles.label, {color: theme.colors.onSurfaceVariant}]}>
                            Add a league to start managing teams, membership, and seasons for this organization.
                        </AppText>
                    </AppCard.Content>
                </AppCard>
            ) : (
                <>
                    <ScrollView>
                        {
                            data.map((league: League) => (
                                <AppCard style={{marginBottom: 8, marginLeft: 1, marginRight: 1}} key={league.id}>
                                    <AppCard.Content style={styles.leagueRow}>
                                        <View style={styles.leagueMeta}>
                                            <AppText style={[styles.title, {color: theme.colors.onSurface}]}>{league.name}</AppText>
                                            <AppText style={[styles.label, {color: theme.colors.onSurfaceVariant}]}>Total Members: {league.totalMembers}</AppText>
                                        </View>
                                        <Link
                                            style={[
                                                styles.link,
                                                {
                                                    backgroundColor: theme.colors.primaryContainer,
                                                    borderColor: theme.colors.primary,
                                                }
                                            ]}
                                            href={`/dashboard/leagues/${league.id}`}
                                        >
                                            <AppText style={[styles.linkText, {color: theme.colors.primary}]}>Details</AppText>
                                        </Link>
                                    </AppCard.Content>
                                </AppCard>
                            ))
                        }
                    </ScrollView>
                    <View style={styles.paginationWrap}>
                        <Pagination
                            currentPage={page}
                            totalItems={total}
                            itemsPerPage={itemsPerPage}
                            onPageChange={setPage}
                        />
                    </View>
                </>
            )}
        </View>
    )
}


const styles = StyleSheet.create({
    link: {
        borderRadius: 8,
        borderWidth: 1,
        paddingHorizontal: 16,
        paddingVertical: 8,
    },
    linkText: {
        fontSize: 14,
        fontWeight: "600",
    },
    title: {
        fontSize: 18,
        fontWeight: "bold",
        marginBottom: 8,
    },
    emptyCard: {
        marginHorizontal: 1,
    },
    emptyContent: {
        gap: 8,
    },
    paginationWrap: {
        alignItems: "flex-end",
        paddingRight: 8,
        paddingTop: 12,
    },
    leagueRow: {
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
        gap: 16,
    },
    leagueMeta: {
        flex: 1,
    },
    label: {
        fontSize: 14,
    },
});

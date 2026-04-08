import {Animated, View, StyleSheet} from "react-native";
import ScrollView = Animated.ScrollView;
import {League, useLeagueData} from "@/hooks/useLeagueData";
import React, {useCallback, useEffect, useState} from "react";
import {useOrganizationStore} from "@/stores/organizationStore";
import Pagination from "@/components/Pagination/Pagination";
import {Link} from "expo-router";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import {useAppTheme} from "@/theme/theme";

interface LeagueListProps {
    refresh?: boolean
}

export function LeagueList(props: LeagueListProps): React.JSX.Element {
    const {fetchData, data, total} = useLeagueData();
    const {organization} = useOrganizationStore()
    const theme = useAppTheme();

    const [page, setPage] = useState(0);
    const itemsPerPage = 5;

    const fetchLeagues = useCallback(async (offset: number) => {
        if (organization !== null) {
            await fetchData({
                organization_id: organization,
                limit: itemsPerPage,
                offset,
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
        fetchLeagues(itemsPerPage * page);
    }, [fetchLeagues, itemsPerPage, page]);

    return (
        <View style={{flex: 1}}>
            <Pagination currentPage={page} totalItems={total} itemsPerPage={itemsPerPage} onPageChange={setPage}/>
                <ScrollView>
                    {
                        data.map((league: League) => (
                            <AppCard style={{marginBottom: 8, marginLeft: 1, marginRight: 1}} key={league.id}>
                                <AppCard.Content>
                                    <AppText style={[styles.title, {color: theme.colors.onSurface}]}>{league.name}</AppText>
                                    <AppText style={[styles.label, {color: theme.colors.onSurfaceVariant}]}>Total Members: {league.totalMembers}</AppText>
                                </AppCard.Content>
                                <AppCard.Actions>
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
                                    </AppCard.Actions>
                            </AppCard>
                        ))
                    }
                </ScrollView>
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
    label: {
        fontSize: 14,
    },
});

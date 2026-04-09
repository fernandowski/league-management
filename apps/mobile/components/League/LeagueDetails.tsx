import {View, StyleSheet} from "react-native";
import React, {useCallback, useEffect} from "react";
import {LeagueDetailResponse, useData} from "@/hooks/useData";
import LeagueDetailsCard from "@/components/League/LeagueDetailsCard";
import MembershipView from "@/components/Membership/MembershipView";
import SeasonTabView from "@/components/Seasons/SeasonTabView";
import Tabs from "@/components/Layout/Tabs";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";

interface LeagueDetailsProps {
    leagueId: string
    refresh: boolean
    selector?: React.ReactNode
}

export default function LeagueDetails(props: LeagueDetailsProps): React.JSX.Element {
    const theme = useAppTheme();

    const {fetchData, data} = useData<LeagueDetailResponse>();

    const fetchLeagueDetails = useCallback(async () => {
        if (!props.leagueId) {
            return;
        }
        await fetchData(`/v1/leagues/${props.leagueId}`);
    }, [fetchData, props.leagueId]);

    const onSeasonAdded = () => {
        fetchLeagueDetails();
    };

    const onSeasonPlanned = () => {
        fetchLeagueDetails();
    };

    useEffect(() => {
        fetchLeagueDetails();
    }, [fetchLeagueDetails, props.refresh]);

    if (!props.leagueId) {
        return (
            <AppCard
                style={[
                    styles.emptyStateCard,
                    {
                        backgroundColor: theme.colors.surfaceVariant,
                        borderColor: theme.colors.outlineVariant,
                    }
                ]}
            >
                <AppCard.Content style={styles.emptyStateContent}>
                    <AppText variant="titleMedium">Choose a league</AppText>
                    <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                        Pick a league from the selector above to load members, season status, and league settings.
                    </AppText>
                </AppCard.Content>
            </AppCard>
        );
    }

    if (!data) {
        return (
            <AppCard
                style={[
                    styles.emptyStateCard,
                    {
                        backgroundColor: theme.colors.surface,
                        borderColor: theme.colors.outline,
                    }
                ]}
            >
                <AppCard.Content style={styles.emptyStateContent}>
                    <AppText variant="titleMedium">Loading league details</AppText>
                    <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                        Fetching the latest league summary and configuration.
                    </AppText>
                </AppCard.Content>
            </AppCard>
        );
    }

    return (
        <View style={styles.container}>
            <AppCard style={[styles.overviewShell, {borderColor: theme.colors.outline}]}>
                <AppCard.Content style={styles.overviewShellContent}>
                    <View style={[styles.selectorSection, {borderBottomColor: theme.colors.outlineVariant}]}>
                        <View style={styles.selectorHeader}>
                            <AppText variant="labelLarge" style={{color: theme.colors.onSurfaceVariant}}>
                                Selected league
                            </AppText>
                        </View>
                        <View style={styles.selectorInput}>
                            {props.selector}
                        </View>
                    </View>
                    <View style={styles.tabSection}>
                        <Tabs tabs={[
                            {key: 'details', title: 'Details', view: <LeagueDetailsCard data={data}/>},
                            {key: 'teams', title: 'Teams', view: <MembershipView leagueId={props.leagueId} onMemberRefresh={fetchLeagueDetails}/>},
                            {key: 'season', title: 'Season', view: <SeasonTabView league={data} onSeasonAdded={onSeasonAdded} onSeasonPlanned={onSeasonPlanned}/>},
                        ]}>
                        </Tabs>
                    </View>
                </AppCard.Content>
            </AppCard>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        gap: 16,
    },
    overviewShell: {
        borderRadius: 18,
    },
    overviewShellContent: {
        gap: 18,
    },
    selectorSection: {
        borderBottomWidth: 1,
        paddingBottom: 16,
        gap: 8,
    },
    selectorHeader: {
        gap: 2,
    },
    selectorInput: {
        width: "100%",
    },
    tabSection: {
        flex: 1
    },
    emptyStateCard: {
        borderRadius: 24,
    },
    emptyStateContent: {
        gap: 8,
        paddingVertical: 8,
    },
});

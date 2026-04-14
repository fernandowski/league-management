import {View, StyleSheet, useWindowDimensions} from "react-native";
import React, {useCallback, useEffect} from "react";
import {SeasonDetailResponse, useData} from "@/hooks/useData";
import SeasonInformation from "@/components/Seasons/SeasonInformation";
import usePost from "@/hooks/usePost";
import { AppButton } from "@/components/ui/AppButton";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import {useAppTheme} from "@/theme/theme";
import SeasonMatchUpManagement from "@/components/Seasons/SeasonMatchUpManagement";
import SeasonStanding from "@/components/Seasons/SeasonStanding";

export interface SeasonDetailsProps {
    seasonId: string
    leagueId: string
    onSeasonPlanned: () => void
}

export default function LeagueSeasonDetails(props: SeasonDetailsProps) {
    const {fetchData, data: seasonDetails} = useData<SeasonDetailResponse>();
    const {postData, result, clearResult} = usePost();
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 768;
    const theme = useAppTheme();
    const showStandings = ['planned', 'in_progress'].indexOf(seasonDetails?.status ?? '') > -1;

    const fetchRoundDetails = useCallback(async () => {
        await fetchData(`/v1/leagues/${props.leagueId}/seasons/${props.seasonId}`);
    }, [fetchData, props.leagueId, props.seasonId]);

    const handlePlanSeason = async () => {
        await postData(`/v1/leagues/${props.leagueId}/seasons/${props.seasonId}/schedules`);
        props.onSeasonPlanned();
        fetchRoundDetails();
    }

    const handleSeasonStart = async () => {
        await postData(`/v1/leagues/${props.leagueId}/seasons/${props.seasonId}/start`);
        props.onSeasonPlanned();
        fetchRoundDetails();
    }

    useEffect(() => {
        clearResult();
        fetchRoundDetails()
    }, [clearResult, fetchRoundDetails]);

    if (!seasonDetails) {
        return (
            <AppCard style={styles.loadingCard}>
                <AppCard.Content>
                    <AppText variant="titleMedium">Loading season</AppText>
                </AppCard.Content>
            </AppCard>
        );
    }

    return (
        <View style={styles.container}>
            {result && (
                <AppCard style={[styles.messageCard, {backgroundColor: theme.colors.errorContainer, borderColor: theme.colors.error}]}>
                    <AppCard.Content>
                        <AppText style={[styles.errorMessage, {color: theme.colors.onErrorContainer}]}>{result}</AppText>
                    </AppCard.Content>
                </AppCard>
            )}

            <View style={[styles.topGrid, !isLargeScreen && styles.topGridStacked]}>
                <AppCard style={[styles.informationCard, {borderColor: theme.colors.outline}]}>
                    <AppCard.Content style={styles.informationCardContent}>
                        <SeasonInformation season={seasonDetails} handleSeasonStart={handleSeasonStart}/>

                        {showStandings && (
                            <View style={[styles.scheduleSection, {borderTopColor: theme.colors.outlineVariant}]}>
                                <View style={styles.matchesHeader}>
                                    <AppText variant="titleMedium">Season schedule</AppText>
                                </View>
                                <SeasonMatchUpManagement seasonId={props.seasonId} onRoundCompleted={fetchRoundDetails}/>
                            </View>
                        )}
                    </AppCard.Content>
                </AppCard>

                {seasonDetails.status === 'pending' && (
                    <AppCard style={[styles.actionCard, {borderColor: theme.colors.outline}]}>
                        <AppCard.Content style={styles.actionCardContent}>
                            <View style={styles.actionTextGroup}>
                                <AppText variant="titleMedium" style={{color: theme.colors.onSurface}}>Plan season</AppText>
                                <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                                    Generate the round-robin schedule once your league members are finalized.
                                </AppText>
                            </View>
                            <View style={styles.planSeasonContainer}>
                                <AppButton variant="submit" onPress={handlePlanSeason}>Plan Season</AppButton>
                            </View>
                        </AppCard.Content>
                    </AppCard>
                )}

                {showStandings && (
                    <SeasonStanding
                        seasonId={props.seasonId}
                        style={styles.standingsCard}
                        contentStyle={styles.standingsCardContent}
                    />
                )}
            </View>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        gap: 16,
    },
    planSeasonContainer: {
        flexDirection: "row",
        justifyContent: "flex-end",
        alignItems: "center",
    },
    topGrid: {
        flexDirection: "row",
        alignItems: "stretch",
        gap: 16,
    },
    topGridStacked: {
        flexDirection: "column",
    },
    informationCard: {
        flex: 1.3,
        borderRadius: 18,
    },
    informationCardContent: {
        gap: 24,
    },
    actionCard: {
        flex: 0.9,
        borderRadius: 18,
    },
    standingsCard: {
        flex: 1,
        minWidth: 0,
    },
    standingsCardContent: {
        flex: 1,
    },
    scheduleSection: {
        borderTopWidth: 1,
        paddingTop: 24,
        gap: 8,
    },
    actionCardContent: {
        flex: 1,
        justifyContent: "space-between",
        gap: 16,
    },
    actionTextGroup: {
        gap: 4,
    },
    matchesHeader: {
        gap: 4,
    },
    loadingCard: {
        borderRadius: 18,
    },
    messageCard: {
        borderRadius: 18,
    },
    errorMessage: {}
});

import {View, StyleSheet, useWindowDimensions} from "react-native";
import React, {useCallback, useEffect} from "react";
import {SeasonDetailResponse, useData} from "@/hooks/useData";
import SeasonInformation from "@/components/Seasons/SeasonInformation";
import usePost from "@/hooks/usePost";
import Matches from "@/components/Seasons/Matches";
import { AppButton } from "@/components/ui/AppButton";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import {useAppTheme} from "@/theme/theme";

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
                    <AppCard.Content>
                        <SeasonInformation season={seasonDetails} handleSeasonStart={handleSeasonStart}/>
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
                                <AppButton onPress={handlePlanSeason}>Plan Season</AppButton>
                            </View>
                        </AppCard.Content>
                    </AppCard>
                )}
            </View>

            {['planned', 'in_progress'].indexOf(seasonDetails.status) > -1 && (
                <View style={styles.matchesContainer}>
                    <Matches season={seasonDetails}/>
                </View>
            )}
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
    actionCard: {
        flex: 0.9,
        borderRadius: 18,
    },
    actionCardContent: {
        flex: 1,
        justifyContent: "space-between",
        gap: 16,
    },
    actionTextGroup: {
        gap: 4,
    },
    matchesContainer: {
        flex: 1,
    },
    loadingCard: {
        borderRadius: 18,
    },
    messageCard: {
        borderRadius: 18,
    },
    errorMessage: {}
});

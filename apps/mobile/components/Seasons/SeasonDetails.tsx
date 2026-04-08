import {Divider, Surface} from "react-native-paper";
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
            <View>
                <AppCard>
                    <AppCard.Title title="Loading..." />
                </AppCard>
            </View>
        );
    }

    return (
        <View>
            <Surface style={{ padding: 18, flexDirection: isLargeScreen ? "row": "column", backgroundColor: theme.colors.surface, borderRadius: 24}}>
                <View>{result && (<AppText style={[styles.errorMessage, {color: theme.colors.error}]}>{result}</AppText>)}</View>
                <Divider/>
                <View style={{flexGrow: 0.5}}>
                    <SeasonInformation season={seasonDetails} handleSeasonStart={handleSeasonStart}/>
                </View>
                <Divider/>
                {
                    seasonDetails.status === 'pending' && (
                        <View>
                            <View style={styles.planSeasonContainer}>
                                <AppText style={{color: theme.colors.onSurface}}>Plan Season</AppText>
                                <AppButton onPress={handlePlanSeason}>Plan Season</AppButton>
                            </View>
                            <Divider/>
                        </View>
                    )
                }
                {
                    ['planned', 'in_progress'].indexOf(seasonDetails.status) > -1
                    &&
                    (
                        <View style={{marginTop: 18, flexGrow: 1}}>
                            <Matches season={seasonDetails}/>
                        </View>
                    )
                }
            </Surface>
        </View>
    )
}

const styles = StyleSheet.create({
    planSeasonContainer: {
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center"
    },
    errorMessage: {}
});

import {Button, Card, Divider, Surface, Text} from "react-native-paper";
import {View, StyleSheet, useWindowDimensions} from "react-native";
import React, {useEffect} from "react";
import {SeasonDetailResponse, useData} from "@/hooks/useData";
import SeasonInformation from "@/components/Seasons/SeasonInformation";
import usePost from "@/hooks/usePost";
import Matches from "@/components/Seasons/Matches";


export interface SeasonDetailsProps {
    seasonId: string
    leagueId: string
    onSeasonPlanned: () => void
}
export default function LeagueSeasonDetails(props: SeasonDetailsProps) {
    const {fetchData, data: seasonDetails, fetching, error} = useData<SeasonDetailResponse>();
    const {postData, result, clearResult} = usePost();
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 768;

    const fetchRoundDetails = async () => {
        await fetchData(`/v1/leagues/${props.leagueId}/seasons/${props.seasonId}`);
    }

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
    }, [props.seasonId, props.leagueId]);

    if (!seasonDetails) {
        return (
            <View>
                <Card>
                    <Card.Title title="Loading..." />
                </Card>
            </View>
        );
    }

    return (
        <View>
            <Surface style={{ padding: 18, flexDirection: isLargeScreen ? "row": "column"}}>
                <View>{result && (<Text style={styles.errorMessage}>{result}</Text>)}</View>
                <Divider/>
                <View style={{flexGrow: 0.5}}>
                    <SeasonInformation season={seasonDetails} handleSeasonStart={handleSeasonStart}/>
                </View>
                <Divider/>
                {
                    seasonDetails.status === 'pending' && (
                        <View>
                            <View style={styles.planSeasonContainer}>
                                <Text>Plan Season</Text>
                                <Button onPress={handlePlanSeason}>Plan Season</Button>
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
    linkText: {
        fontSize: 14,
        color: 'rgb(103, 80, 164)'
    },
    planSeasonContainer: {
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center"
    },
    errorMessage: {
        color: "red"
    }
});



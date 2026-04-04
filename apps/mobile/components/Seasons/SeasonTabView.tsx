import {StyleSheet, View} from "react-native";
import React from "react";
import {LeagueDetailResponse} from "@/hooks/useData";
import LeagueSeasonDetails from "@/components/Seasons/SeasonDetails";
import AddSeason from "@/components/Seasons/AddSeason";

interface SeasonsViewProps {
    league: LeagueDetailResponse
    onSeasonAdded: () =>  void
    onSeasonPlanned: () => void
}

const SeasonTabView = (props: SeasonsViewProps) => {
    return (
        <View style={styles.container}>
            {
                props.league.season ?
                (<LeagueSeasonDetails seasonId={props.league.season.id} leagueId={props.league.id} onSeasonPlanned={props.onSeasonPlanned}/>)
                :
                (<AddSeason leagueId={props.league.id} onSeasonAdded={props.onSeasonAdded}/>)
            }
        </View>
    )
}

export default SeasonTabView

const styles = StyleSheet.create({
    container: {
        flex: 1,
    },
    buttonContainer: {
        marginVertical: 27,
        alignItems: "center",
        gap: 8
    },
    surfaceContainer: {
        marginTop: 27,
        gap: 18,
        height: "80%"
    }
})

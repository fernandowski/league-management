import {MatchesResponse, MatchScore} from "@/hooks/useData";
import {View, StyleSheet} from "react-native";
import {Divider, Text} from "react-native-paper";
import MatchUpScoreCard from "@/components/Seasons/MatchUpScoreCard";

export interface SeasonMatchUpsListProps {
    data: MatchesResponse;
    onMatchPress: (matchScore: MatchScore) => void;
}

export default function MatchUpRound(props: SeasonMatchUpsListProps) {
    return (
        <View>
            <View><Text style={styles.roundText}>Round {props.data.round}</Text></View>
            <Divider/>
            <View style={[styles.matchUpCardContainer]}>
                {
                    props.data.matches.map((match: MatchScore) => (
                        <View style={[styles.matchUpScoreCardContainer]} key={match.id}>
                            <MatchUpScoreCard data={match} onPress={props.onMatchPress}/>
                        </View>
                        ))
                }
            </View>
        </View>
    )
}

const styles = StyleSheet.create({
    matchUpCardContainer: {
        justifyContent: "center",
        gap: 16,
        flexWrap: "wrap",
        marginTop: 16
    },
    matchUpScoreCardContainer: {
    },
    roundText: {
        fontSize: 18,
        fontWeight: "bold",
        marginBottom: 8,
    },
})

import {MatchesResponse, MatchScore} from "@/hooks/useData";
import {View, StyleSheet} from "react-native";
import {Divider} from "react-native-paper";
import MatchUpScoreCard from "@/components/Seasons/MatchUpScoreCard";
import { AppText } from "@/components/ui/AppText";

export interface SeasonMatchUpsListProps {
    data: MatchesResponse;
    onMatchPress: (matchScore: MatchScore) => void;
    isEditableRound?: boolean;
}

export default function MatchUpRound(props: SeasonMatchUpsListProps) {
    return (
        <View>
            <View><AppText style={styles.roundText}>Round {props.data.round}</AppText></View>
            <Divider/>
            <View style={[styles.matchUpCardContainer]}>
                {
                    props.data.matches.map((match: MatchScore) => (
                        <View style={[styles.matchUpScoreCardContainer]} key={match.id}>
                            <MatchUpScoreCard
                                data={match}
                                onPress={props.onMatchPress}
                                disabled={!props.isEditableRound || match.home_team === "Bye" || match.away_team === "Bye"}
                            />
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

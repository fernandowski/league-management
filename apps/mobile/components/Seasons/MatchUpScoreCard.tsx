import {MatchScore} from "@/hooks/useData";
import {Surface, Text} from "react-native-paper";
import {TouchableOpacity, View, StyleSheet, useWindowDimensions} from "react-native";

export interface MatchUpScoreCardProps {
    data: MatchScore;
    onPress: (matchScore: MatchScore) => void
}

export default function MatchUpScoreCard(props: MatchUpScoreCardProps) {
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 768;
    return (
        <Surface style={{flex: 1}}>
            <TouchableOpacity style={{paddingVertical: 18}} onPress={() => props.onPress(props.data)}>
                <View style={[styles.row, styles.center]}>
                    <View style={[styles.row]}>
                        <Text style={[
                             {textAlign: "right"},
                            isLargeScreen ? styles.labelWidth : styles.smallScreenLabelWidth
                        ]}>{props.data.home_team}</Text>
                    </View>
                    <View style={[styles.score, styles.row]}>
                        <Text>{props.data.home_score}</Text>
                        <Text> - </Text>
                        <Text>{props.data.away_score}</Text>
                    </View>
                    <View style={[styles.row,]}>
                        <Text style={[isLargeScreen ? styles.labelWidth : styles.smallScreenLabelWidth,{textAlign: "left"} ]}>{props.data.away_team}</Text>
                    </View>
                </View>
            </TouchableOpacity>
        </Surface>
    )
}

const styles = StyleSheet.create({
    teamContainer: {
        flex: 1,
        flexDirection: "row",
        alignItems: "center",
    },
    labelWidth: {
        width: 300
    },
    smallScreenLabelWidth: {
        width: 100
    },
    score: {
        marginVertical: 4,
        marginHorizontal: 16,
    },
    row: {
        flexDirection: "row",
    },
    center: {
        justifyContent: "center",
        alignItems: "center"
    }
});

import {MatchScore} from "@/hooks/useData";
import {Surface} from "react-native-paper";
import {Pressable, View, StyleSheet, useWindowDimensions} from "react-native";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";

export interface MatchUpScoreCardProps {
    data: MatchScore;
    onPress: (matchScore: MatchScore) => void
    disabled?: boolean
}

export default function MatchUpScoreCard(props: MatchUpScoreCardProps) {
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 768;
    const theme = useAppTheme();
    return (
        <Surface style={[{flex: 1}, props.disabled && styles.disabledSurface]}>
            <Pressable
                style={({ pressed }) => [
                    styles.pressable,
                    props.disabled && styles.disabledPressable,
                    !props.disabled && pressed && {backgroundColor: theme.colors.primaryContainer}
                ]}
                disabled={props.disabled}
                onPress={() => props.onPress(props.data)}
            >
                <View style={[styles.row, styles.center]}>
                    <View style={[styles.row]}>
                        <AppText style={[
                             {textAlign: "right"},
                            isLargeScreen ? styles.labelWidth : styles.smallScreenLabelWidth
                        ]}>{props.data.home_team}</AppText>
                    </View>
                    <View style={[styles.score, styles.row]}>
                        <AppText>{props.data.home_score}</AppText>
                        <AppText> - </AppText>
                        <AppText>{props.data.away_score}</AppText>
                    </View>
                    <View style={[styles.row,]}>
                        <AppText style={[isLargeScreen ? styles.labelWidth : styles.smallScreenLabelWidth,{textAlign: "left"} ]}>{props.data.away_team}</AppText>
                    </View>
                </View>
            </Pressable>
        </Surface>
    )
}

const styles = StyleSheet.create({
    pressable: {
        paddingVertical: 18,
    },
    disabledSurface: {
        opacity: 0.6,
    },
    disabledPressable: {
        cursor: "default",
    },
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

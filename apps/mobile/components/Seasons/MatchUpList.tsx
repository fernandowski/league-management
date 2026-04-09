import {Match, Round} from "@/hooks/useData";
import {StyleSheet, View} from "react-native";
import {Divider} from "react-native-paper";
import MatchUp from "@/components/Seasons/MatchUp";
import { AppText } from "@/components/ui/AppText";
import {useAppTheme} from "@/theme/theme";

export interface MatchUpListProps {
    round: Round
}

export default function MatchUpList(props: MatchUpListProps ) {
    const theme = useAppTheme();
    return (
        <View style={styles.container}>
            <View style={[styles.header, {backgroundColor: theme.colors.surfaceVariant}]}>
                <View style={styles.headerContent}>
                    <AppText variant="titleMedium" style={{color: theme.colors.onSurface}}>Round {props.round.round_number}</AppText>
                    <AppText variant="bodySmall" style={{color: theme.colors.onSurfaceVariant}}>
                        {props.round.matches.length} match{props.round.matches.length === 1 ? "" : "es"}
                    </AppText>
                </View>
            </View>
            <View style={styles.matches}>
                {
                    props.round.matches.map((match: Match) => {
                        return (
                            <View key={match.id}>
                                <MatchUp match={match}/>
                                <Divider/>
                            </View>
                        )
                    })
                }
            </View>
        </View>
    )
}


const styles = StyleSheet.create({
    container: {
        borderRadius: 20,
        overflow: "hidden",
    },
    header: {
        justifyContent: "center",
        minHeight: 56,
    },
    headerContent: {
        paddingHorizontal: 18,
        paddingVertical: 12,
        gap: 2,
    },
    matches: {
        gap: 4,
        paddingTop: 4,
    },
})

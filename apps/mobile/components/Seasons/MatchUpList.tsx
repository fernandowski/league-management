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
        <View>
            <View style={[styles.header, {backgroundColor: theme.colors.surfaceVariant}]}>
                <View style={{paddingLeft: 18}}><AppText style={{color: theme.colors.onSurface}}>MatchUps For Round {props.round.round_number }</AppText></View>
            </View>
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
    )
}


const styles = StyleSheet.create({
    header: {
        justifyContent: "center",
        borderTopLeftRadius: 16,
        borderTopRightRadius: 16,
        minHeight: 44,
    }
})

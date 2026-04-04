import {TouchableOpacity, View, StyleSheet} from "react-native";
import {IconButton, Text} from "react-native-paper";
import {container} from "ansi-fragments";

export interface MatchesPaginationProps {
    currentRound: string
    onPrevious: () => void
    onNext: () => void
}
export default function MatchesPagination(props: MatchesPaginationProps) {
    return (
        <View style={styles.container}>
            <TouchableOpacity>
                <IconButton icon="chevron-left" size={20} onPress={props.onPrevious}/>
            </TouchableOpacity>
            <View><Text style={styles.roundText}>Round {props.currentRound}</Text></View>
            <TouchableOpacity>
                <IconButton icon="chevron-right" size={20} onPress={props.onNext}/>
            </TouchableOpacity>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        flexDirection: "row",
        alignItems: "center",
    },
    roundText: {
        fontWeight: "bold"
    }
});

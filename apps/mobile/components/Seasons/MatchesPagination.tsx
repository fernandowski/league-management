import {View, StyleSheet} from "react-native";
import {IconButton} from "react-native-paper";

import { AppText } from "@/components/ui/AppText";

export interface MatchesPaginationProps {
    currentRound: string
    onPrevious: () => void
    onNext: () => void
    disablePrevious?: boolean
    disableNext?: boolean
}
export default function MatchesPagination(props: MatchesPaginationProps) {
    return (
        <View style={styles.container}>
            <IconButton icon="chevron-left" size={20} onPress={props.onPrevious} disabled={props.disablePrevious}/>
            <View><AppText style={styles.roundText}>Round {props.currentRound}</AppText></View>
            <IconButton icon="chevron-right" size={20} onPress={props.onNext} disabled={props.disableNext}/>
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

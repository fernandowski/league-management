import {Match} from "@/hooks/useData";
import {View, StyleSheet, Pressable} from "react-native";
import { AppText } from "@/components/ui/AppText";
import {useAppTheme} from "@/theme/theme";


export interface MatchUpProps {
    match: Match
}

export default function MatchUp(props: MatchUpProps) {
    const theme = useAppTheme();
    return (
        <Pressable style={({pressed}) => [
            styles.container,
            pressed && {backgroundColor: theme.colors.primaryContainer}
        ]}>
            <View>
                <AppText style={[styles.label, {color: theme.colors.onSurface}]}>{props.match.home_team}</AppText>
            </View>
            <View>
                <AppText style={{color: theme.colors.onSurfaceVariant}}>VS</AppText>
            </View>
            <View>
                <AppText style={[styles.label, {color: theme.colors.onSurface}]}>{props.match.away_team}</AppText>
            </View>
        </Pressable>
    )
}


const styles = StyleSheet.create({
    container: {
        flexDirection: "row",
        justifyContent: "space-evenly",
        paddingVertical: 18,
        borderRadius: 16,
    },
    label: {
        fontWeight: "bold",
        width: 80,
    },
})

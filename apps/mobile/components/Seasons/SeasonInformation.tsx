import {SeasonDetailResponse} from "@/hooks/useData";
import {StyleSheet, View} from "react-native";
import { AppButton } from "@/components/ui/AppButton";
import { AppText } from "@/components/ui/AppText";

export interface SeasonInformationProps {
    season: SeasonDetailResponse,
    handleSeasonStart: () => void
}
export default function SeasonInformation(props: SeasonInformationProps) {
    return (
            <View style={styles.surfaceCard}>
                    <View style={{flexDirection: "row", justifyContent: "space-between"}}>
                        <View style={[styles.row]}>
                            <AppText style={styles.label}>Name: </AppText>
                            <AppText style={styles.value}>{props.season.name} </AppText>
                        </View>
                        {
                            ['planned'].indexOf(props.season.status) > -1 && (
                                <View>
                                    <AppButton onPress={props.handleSeasonStart}>Start Season</AppButton>
                                </View>
                            )
                        }
                    </View>
                    <View style={styles.row}>
                        <AppText style={styles.label}>Status: </AppText>
                        <AppText style={styles.value}>{props.season.status} </AppText>
                    </View>
                    {
                        props.season.rounds.length > 0 ?? (
                            <View style={styles.row}>
                                <AppText style={styles.label}>Status: </AppText>
                                <AppText style={styles.value}>{props.season.status} </AppText>
                            </View>
                        )
                    }
            </View>
    )
}


const styles = StyleSheet.create({
    surfaceCard: {
        padding: 18
    },
    row: {
        marginBottom: 8,
        flexGrow: 0,
        flexDirection: "row",
        alignItems: "center"
    },
    label: {
        fontWeight: "bold",
        paddingRight: 18,
        width: 80
    },
    value: {
        flex: 1
    }
})

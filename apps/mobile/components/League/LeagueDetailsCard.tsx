import {View, StyleSheet} from "react-native";
import React from "react";
import {LeagueDetailResponse} from "@/hooks/useData";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";

export interface LeagueDetailsCardProps {
    data: LeagueDetailResponse
}
const LeagueDetailsCard = (props: LeagueDetailsCardProps) => {
    return (
        <AppCard>
            <AppCard.Title title={"League Info"} titleStyle={{fontWeight: "bold"}}></AppCard.Title>
            <AppCard.Content>
                <View style={{flexDirection: "row",  flexWrap: "wrap", gap: 18}}>
                    <View>
                        <View style={styles.row}>
                            <AppText style={styles.label}>League ID: </AppText>
                            <AppText>{props.data.id}</AppText>
                        </View>
                        <View style={styles.row}>
                            <AppText style={styles.label}>Name: </AppText>
                            <AppText>{props.data.name}</AppText>
                        </View>
                        <View style={styles.row}>
                            <AppText style={styles.label}>Active Members: </AppText>
                            <AppText>{props.data.active_members}</AppText>
                        </View>
                    </View>
                    <View style={{width: "100%"}}>
                        {
                            props.data.season ?
                                (
                                    <View>
                                        <View style={styles.row}>
                                            <AppText style={styles.label}>Season ID: </AppText>
                                            <AppText>{props.data.season.id}</AppText>
                                        </View>
                                        <View style={styles.row}>
                                            <AppText style={styles.label}>Name: </AppText>
                                            <AppText>{props.data.season.name}</AppText>
                                        </View>
                                        <View style={styles.row}>
                                            <AppText style={styles.label}>Status: </AppText>
                                            <AppText>{props.data.season.status}</AppText>
                                        </View>
                                    </View>
                                )
                                    :
                                    (
                                        <View>
                                            <AppText style={styles.messageText}> No Season has been configured. Go to Seasons section to configure one. </AppText>
                                        </View>
                                    )
                        }

                    </View>
                </View>
            </AppCard.Content>
        </AppCard>
    )
}
export default LeagueDetailsCard;


const styles = StyleSheet.create({
    row: {
        marginBottom: 8,
        alignItems: "flex-start",
    },
    label: {
        fontWeight: "bold",
    },
    messageText: {
        textAlign: "left"
    }
})

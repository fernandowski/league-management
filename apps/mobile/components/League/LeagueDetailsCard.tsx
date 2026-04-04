import {Card, Text} from "react-native-paper";
import {View, StyleSheet} from "react-native";
import React from "react";
import {LeagueDetailResponse} from "@/hooks/useData";

export interface LeagueDetailsCardProps {
    data: LeagueDetailResponse
}
const LeagueDetailsCard = (props: LeagueDetailsCardProps) => {
    return (
        <Card>
            <Card.Title title={"League Info"} titleStyle={{fontWeight: "bold"}}></Card.Title>
            <Card.Content>
                <View style={{flexDirection: "row",  flexWrap: "wrap", gap: 18}}>
                    <View>
                        <View style={styles.row}>
                            <Text style={styles.label}>League ID: </Text>
                            <Text>{props.data.id}</Text>
                        </View>
                        <View style={styles.row}>
                            <Text style={styles.label}>Name: </Text>
                            <Text>{props.data.name}</Text>
                        </View>
                        <View style={styles.row}>
                            <Text style={styles.label}>Active Members: </Text>
                            <Text>{props.data.active_members}</Text>
                        </View>
                    </View>
                    <View style={{width: "100%"}}>
                        {
                            props.data.season ?
                                (
                                    <View>
                                        <View style={styles.row}>
                                            <Text style={styles.label}>Season ID: </Text>
                                            <Text>{props.data.season.id}</Text>
                                        </View>
                                        <View style={styles.row}>
                                            <Text style={styles.label}>Name: </Text>
                                            <Text>{props.data.season.name}</Text>
                                        </View>
                                        <View style={styles.row}>
                                            <Text style={styles.label}>Status: </Text>
                                            <Text>{props.data.season.status}</Text>
                                        </View>
                                    </View>
                                )
                                    :
                                    (
                                        <View>
                                            <Text style={styles.messageText}> No Season has been configured. Go to Seasons section to configure one. </Text>
                                        </View>
                                    )
                        }

                    </View>
                </View>
            </Card.Content>
        </Card>
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

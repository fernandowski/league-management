import {Match} from "@/hooks/useData";
import {View, StyleSheet, Pressable} from "react-native";
import {Text} from "react-native-paper";
import {useState} from "react";


export interface MatchUpProps {
    match: Match
}

export default function MatchUp(props: MatchUpProps) {
    return (
        <Pressable style={styles.container}>
            <View>
                <Text style={styles.label}>{props.match.home_team}</Text>
            </View>
            <View>
                <Text>VS</Text>
            </View>
            <View>
                <Text style={styles.label}>{props.match.away_team}</Text>
            </View>
        </Pressable>
    )
}


const styles = StyleSheet.create({
    container: {
        flexDirection: "row",
        justifyContent: "space-evenly",
        paddingVertical: 18
    },
    label: {
        fontWeight: "bold",
        width: 80,
    },
    onHover: {
        backgroundColor: "green"
    }
})

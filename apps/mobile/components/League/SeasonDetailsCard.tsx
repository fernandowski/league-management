import {LeagueDetailsCardProps} from "@/components/League/LeagueDetailsCard";
import {Card, Text} from "react-native-paper";
import {StyleSheet, View} from "react-native";
import React from "react";

export default function SeasonDetailsCard(props: LeagueDetailsCardProps) {
    return (
        <Card>
            <Card.Title title={"Current Season"}></Card.Title>
            <Card.Content>
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
            </Card.Content>
        </Card>
    )
}


const styles = StyleSheet.create({
    row: {
        marginBottom: 8,
        alignItems: "flex-start",
    },
    label: {
        fontWeight: "bold"
    }
})

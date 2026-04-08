import {LeagueDetailsCardProps} from "@/components/League/LeagueDetailsCard";
import {StyleSheet, View} from "react-native";
import React from "react";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";

export default function SeasonDetailsCard(props: LeagueDetailsCardProps) {
    return (
        <AppCard>
            <AppCard.Title title={"Current Season"}></AppCard.Title>
            <AppCard.Content>
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
            </AppCard.Content>
        </AppCard>
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

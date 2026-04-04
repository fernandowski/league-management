import {Match, Round} from "@/hooks/useData";
import {Pressable, StyleSheet, View} from "react-native";
import {Divider, Text} from "react-native-paper";
import MatchUp from "@/components/Seasons/MatchUp";
import {useRef, useState} from "react";

export interface MatchUpListProps {
    round: Round
}

export default function MatchUpList(props: MatchUpListProps ) {
    return (
        <View>
            <View style={styles.header}>
                <View style={{paddingLeft: 18}}><Text>MatchUps For Round {props.round.round_number }</Text></View>
            </View>
                {
                    props.round.matches.map((match: Match) => {
                        return (
                            <View>
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
        backgroundColor: "#c7cdd2",
        justifyContent: "center"
    },
    onHover: {
        backgroundColor: "purple"
    }
})

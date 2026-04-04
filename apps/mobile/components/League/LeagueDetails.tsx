import {View, StyleSheet, useWindowDimensions} from "react-native";
import React, {useEffect, useState} from "react";
import {LeagueDetailResponse, useData} from "@/hooks/useData";
import LeagueDetailsCard from "@/components/League/LeagueDetailsCard";
import MembershipView from "@/components/Membership/MembershipView";
import SeasonTabView from "@/components/Seasons/SeasonTabView";
import Tabs from "@/components/Layout/Tabs";

interface LeagueDetailsProps {
    leagueId: string
    refresh: boolean
}


export default function LeagueDetails(props: LeagueDetailsProps): React.JSX.Element {
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 768;

    const {fetchData, data} = useData<LeagueDetailResponse>();

    const fetch = async () => {
        fetchData(`/v1/leagues/${props.leagueId}`);
    }

    const onMemberRefresh = () => {
        fetch();
    }

    const onSeasonAdded = () => {
        fetch();
    }

    const onSeasonPlanned = () => {
        fetch();
    }

    useEffect(() => {
        fetch();
    }, [props.leagueId, props.refresh]);

    if (!data) {
        return <></>
    }

    return (
        <View style={styles.container}>
            <View style={ isLargeScreen ? styles.rowContainer : styles.columnContainer}>
                <View style={isLargeScreen ? {flex: 0.5} : styles.fullWidthCard}>
                    <LeagueDetailsCard data={data}/>
                </View>
            </View>
            <View style={styles.tabsContainer}>
                <Tabs tabs={[
                    {key: 'membership', title: 'Membership', view: <MembershipView leagueId={props.leagueId} onMemberRefresh={onMemberRefresh}/>},
                    {key: 'season', title: 'Season', view: <SeasonTabView league={data} onSeasonAdded={onSeasonAdded} onSeasonPlanned={onSeasonPlanned}/>},
                ]}>
                </Tabs>
            </View>
        </View>
    )
}




const styles = StyleSheet.create({
    container: {
        flex: 1,
    },
    rowContainer: {
        flexDirection: "row",
        gap: 18,
        paddingHorizontal: 16,
        paddingVertical: 8,
    },
    columnContainer: {
        flexDirection: "column",
        alignItems: "center",
        gap: 12,
        paddingHorizontal: 16,
        paddingVertical: 8,
    },
    fullWidthCard: {
        width: "100%",
    },
    tabsContainer: {
        flex: 1
    },
});

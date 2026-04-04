import {Animated, View, StyleSheet} from "react-native";
import ScrollView = Animated.ScrollView;
import {Button, Card, Text} from "react-native-paper";
import {League, useLeagueData} from "@/hooks/useLeagueData";
import React, {useEffect, useState} from "react";
import {useOrganizationStore} from "@/stores/organizationStore";
import Pagination from "@/components/Pagination/Pagination";
import {Link} from "expo-router";

interface LeagueListProps {
    refresh?: boolean
}

export function LeagueList(props: LeagueListProps): React.JSX.Element {
    const {fetchData, fetching, error, data, total} = useLeagueData();
    const {organization} = useOrganizationStore()

    const [page, setPage] = useState(0);
    const [itemsPerPage, setItemsPerPage] = useState(5);

    useEffect(() => {
        if (organization !== null) {
            setPage(0);
            fetchData({organization_id: organization, limit: itemsPerPage, offset: 0, term: ""});
        }

    }, [organization, props.refresh]);


    useEffect(() => {
        if (organization !== null) {
            fetchData({
                organization_id: organization,
                limit: itemsPerPage,
                offset: itemsPerPage * page,
                term: "",
            });
        }
    }, [page, itemsPerPage]);

    return (
        <View style={{flex: 1}}>
            <Pagination currentPage={page} totalItems={total} itemsPerPage={itemsPerPage} onPageChange={setPage}/>
                <ScrollView>
                    {
                        data.map((league: League) => (
                            <Card style={{marginBottom: 8, marginLeft: 1, marginRight: 1}} key={league.id}>
                                <Card.Content>
                                    <Text style={styles.title}>{league.name}</Text>
                                    <Text style={styles.label}>Total Members: {league.totalMembers}</Text>
                                </Card.Content>
                                <Card.Actions>
                                    <Link style={styles.link} href={`/dashboard/leagues/${league.id}`}><Text style={styles.linkText}>Details</Text></Link>
                                </Card.Actions>
                            </Card>
                        ))
                    }
                </ScrollView>
        </View>
    )
}


const styles = StyleSheet.create({
    link: {
        borderRadius: 8,
        paddingHorizontal: 16,
        paddingVertical: 8,
        backgroundColor: 'rgba(0, 0, 0, 0.00)',
        borderColor: "#0056b3"
    },
    linkText: {
        fontSize: 14,
        color: 'rgb(103, 80, 164)'
    },
    title: {
        fontSize: 18,
        fontWeight: "bold",
        marginBottom: 8,
    },
    label: {
        fontSize: 14,
        color: "#555",
    },
});

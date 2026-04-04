import {StyleSheet, TouchableOpacity, View} from "react-native";
import {Card, DataTable, Text} from "react-native-paper";
import React, {useEffect, useState} from "react";
import {useData} from "@/hooks/useData";
import {useOrganizationStore} from "@/stores/organizationStore";
import Pagination from "@/components/Pagination/Pagination";
import {useRouter} from "expo-router";

interface Season {
    id: string;
    name: string;
    league_id: string;
    status: string
}

interface SeasonResponse {
    data: Season[];
    total: number;
}

interface SeasonTableIProps {
    leagueId: string;
}

export default function SeasonsTable(props: SeasonTableIProps) {
    const {organization} = useOrganizationStore();
    const router = useRouter()

    const [page, setPage] = useState(0);
    const [seasons, setSeasons] = useState<Season[]>([]);
    const [total, setTotal] = useState(0);
    const [numberOfItemsPerPage, onItemsPerPageChange] = useState(10);

    const {fetchData} = useData<SeasonResponse>();

    const from = page * numberOfItemsPerPage;

    useEffect(() => {
        const fetchSeasons = async () => {
            const offsetFilter = `offset=${encodeURIComponent(from)}`;
            const limitFilter = `limit=${encodeURIComponent(numberOfItemsPerPage)}`;
            const response = await fetchData(
                `/v1/leagues/${props.leagueId}/seasons?${[offsetFilter, limitFilter].join(
                    "&"
                )}`
            );
            if (response) {
                setSeasons(response.data);
                setTotal(response.total);
            }
        };

        if (props.leagueId) {
            fetchSeasons();
        }
    }, [props.leagueId, page, numberOfItemsPerPage]);

    useEffect(() => {
        if (page !== 0) {
            setPage(0);
        }
    }, [numberOfItemsPerPage, organization]);

    return (
        <View>
            <Pagination currentPage={page} totalItems={total} itemsPerPage={10} onPageChange={setPage}/>
            {
                seasons.map((season: Season) => (
                    <Card style={{marginBottom: 8, marginLeft: 1, marginRight: 1}} key={season.id}>
                        <Card.Content>
                            <View>
                                <View>
                                    <Text style={styles.title}>{season.name}</Text>
                                </View>
                                <View>
                                    <Text style={styles.label}>
                                        Status: <Text>{season.status}</Text>
                                    </Text>
                                </View>
                            </View>
                        </Card.Content>
                        <Card.Actions>
                            <TouchableOpacity
                                style={styles.link}
                                onPress={() => router.push(`/dashboard/seasons/${season.id}`)}>
                                <Text style={styles.linkText}>Details</Text>
                            </TouchableOpacity>
                        </Card.Actions>
                    </Card>
                ))
            }
        </View>
    );
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
    paginationContainer: {
        alignSelf: "flex-end"
    }
});

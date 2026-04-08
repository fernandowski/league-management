import {StyleSheet, View} from "react-native";
import React, {useCallback, useEffect, useState} from "react";
import {useData} from "@/hooks/useData";
import {useOrganizationStore} from "@/stores/organizationStore";
import Pagination from "@/components/Pagination/Pagination";
import {useRouter} from "expo-router";
import { AppButton } from "@/components/ui/AppButton";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import {useAppTheme} from "@/theme/theme";

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
    const theme = useAppTheme();

    const [page, setPage] = useState(0);
    const [seasons, setSeasons] = useState<Season[]>([]);
    const [total, setTotal] = useState(0);
    const [numberOfItemsPerPage] = useState(10);

    const {fetchData} = useData<SeasonResponse>();

    const from = page * numberOfItemsPerPage;

    const fetchSeasons = useCallback(async () => {
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
    }, [fetchData, from, numberOfItemsPerPage, props.leagueId]);

    useEffect(() => {
        if (props.leagueId) {
            fetchSeasons();
        }
    }, [fetchSeasons, props.leagueId]);

    useEffect(() => {
        if (page !== 0) {
            setPage(0);
        }
    }, [numberOfItemsPerPage, organization, page]);

    return (
        <View>
            <Pagination currentPage={page} totalItems={total} itemsPerPage={10} onPageChange={setPage}/>
            {
                seasons.map((season: Season) => (
                    <AppCard style={{marginBottom: 8, marginLeft: 1, marginRight: 1}} key={season.id}>
                        <AppCard.Content>
                            <View>
                                <View>
                                    <AppText style={[styles.title, {color: theme.colors.onSurface}]}>{season.name}</AppText>
                                </View>
                                <View>
                                    <AppText style={[styles.label, {color: theme.colors.onSurfaceVariant}]}>
                                        Status: {season.status}
                                    </AppText>
                                </View>
                            </View>
                        </AppCard.Content>
                        <AppCard.Actions>
                            <AppButton
                                mode="contained-tonal"
                                style={[
                                    styles.link,
                                    {
                                        backgroundColor: theme.colors.primaryContainer,
                                        borderColor: theme.colors.primary,
                                    }
                                ]}
                                onPress={() => router.push(`/dashboard/seasons/${season.id}`)}>
                                <AppText style={[styles.linkText, {color: theme.colors.primary}]}>Details</AppText>
                            </AppButton>
                        </AppCard.Actions>
                    </AppCard>
                ))
            }
        </View>
    );
}


const styles = StyleSheet.create({
    link: {
        borderRadius: 8,
        borderWidth: 1,
        paddingHorizontal: 16,
        paddingVertical: 8,
    },
    linkText: {
        fontSize: 14,
        fontWeight: "600",
    },
    title: {
        fontSize: 18,
        fontWeight: "bold",
        marginBottom: 8,
    },
    label: {
        fontSize: 14,
    },
    paginationContainer: {
        alignSelf: "flex-end"
    }
});

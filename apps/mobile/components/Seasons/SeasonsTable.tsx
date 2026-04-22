import {StyleSheet, View} from "react-native";
import React, {useCallback, useEffect, useState} from "react";
import {SeasonSearchResponse, SeasonSummaryResponse, useData} from "@/hooks/useData";
import {useOrganizationStore} from "@/stores/organizationStore";
import Pagination from "@/components/Pagination/Pagination";
import {useRouter} from "expo-router";
import { AppButton } from "@/components/ui/AppButton";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import {useAppTheme} from "@/theme/theme";

interface SeasonTableProps {
    leagueId: string;
    statusFilter?: string;
    phaseFilter?: string;
    emptyMessage?: string;
}

export default function SeasonsTable(props: SeasonTableProps) {
    const {organization} = useOrganizationStore();
    const router = useRouter();
    const theme = useAppTheme();

    const [page, setPage] = useState(0);
    const [seasons, setSeasons] = useState<SeasonSummaryResponse[]>([]);
    const [total, setTotal] = useState(0);
    const [numberOfItemsPerPage] = useState(10);

    const {fetchData} = useData<SeasonSearchResponse>();

    const from = page * numberOfItemsPerPage;

    const fetchSeasons = useCallback(async () => {
        if (!props.leagueId) {
            setSeasons([]);
            setTotal(0);
            return;
        }

        const offsetFilter = `offset=${encodeURIComponent(from)}`;
        const limitFilter = `limit=${encodeURIComponent(numberOfItemsPerPage)}`;
        const statusFilter = props.statusFilter ? `status=${encodeURIComponent(props.statusFilter)}` : "";
        const phaseFilter = props.phaseFilter ? `phase=${encodeURIComponent(props.phaseFilter)}` : "";
        const response = await fetchData(
            `/v1/leagues/${props.leagueId}/seasons?${[offsetFilter, limitFilter, statusFilter, phaseFilter].filter(Boolean).join("&")}`
        );
        if (response) {
            setSeasons(response.data);
            setTotal(response.total);
        }
    }, [fetchData, from, numberOfItemsPerPage, props.leagueId, props.phaseFilter, props.statusFilter]);

    useEffect(() => {
        if (props.leagueId) {
            fetchSeasons();
        }
    }, [fetchSeasons, props.leagueId]);

    useEffect(() => {
        if (page !== 0) {
            setPage(0);
        }
    }, [numberOfItemsPerPage, organization, page, props.leagueId, props.phaseFilter, props.statusFilter]);

    return (
        <View style={styles.container}>
            <Pagination currentPage={page} totalItems={total} itemsPerPage={10} onPageChange={setPage}/>
            {seasons.length === 0 ? (
                <View style={[styles.emptyState, {borderColor: theme.colors.outlineVariant, backgroundColor: theme.colors.surfaceVariant}]}>
                    <AppText variant="titleSmall">No seasons found</AppText>
                    <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                        {props.emptyMessage ?? "No seasons match this view yet."}
                    </AppText>
                </View>
            ) : null}
            {seasons.map((season) => (
                <AppCard style={styles.card} key={season.id}>
                    <AppCard.Content style={styles.cardContent}>
                        <View style={styles.cardHeader}>
                            <View style={styles.titleGroup}>
                                <AppText style={[styles.title, {color: theme.colors.onSurface}]}>{season.name}</AppText>
                                <AppText style={[styles.label, {color: theme.colors.onSurfaceVariant}]}>
                                    {formatStatus(season.status)} · {formatLabel(season.phase)}
                                </AppText>
                            </View>
                            <View style={[styles.statusPill, {backgroundColor: theme.colors.surfaceVariant}]}>
                                <AppText variant="labelMedium" style={{color: theme.colors.onSurfaceVariant}}>
                                    {formatStatus(season.status)}
                                </AppText>
                            </View>
                        </View>
                        <View style={styles.metaGrid}>
                            <MetaItem label="Phase" value={formatLabel(season.phase)} />
                            <MetaItem label="Champion" value={season.champion_team_name ?? "Not decided"} />
                        </View>
                    </AppCard.Content>
                    <AppCard.Actions>
                        <AppButton variant="secondary" onPress={() => router.push(`/dashboard/seasons/${season.id}`)}>
                            View details
                        </AppButton>
                    </AppCard.Actions>
                </AppCard>
            ))}
        </View>
    );
}

function MetaItem({label, value}: {label: string; value: string}) {
    const theme = useAppTheme();

    return (
        <View style={styles.metaItem}>
            <AppText variant="labelMedium" style={{color: theme.colors.onSurfaceVariant}}>
                {label}
            </AppText>
            <AppText variant="bodyMedium">
                {value}
            </AppText>
        </View>
    );
}

function formatStatus(value: string) {
    return formatLabel(value);
}

function formatLabel(value: string) {
    return value.replaceAll("_", " ").replace(/\b\w/g, (char) => char.toUpperCase());
}

const styles = StyleSheet.create({
    container: {
        gap: 12,
    },
    card: {
        marginBottom: 8,
        marginLeft: 1,
        marginRight: 1,
    },
    cardContent: {
        gap: 16,
    },
    cardHeader: {
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "flex-start",
        gap: 12,
        flexWrap: "wrap",
    },
    titleGroup: {
        flex: 1,
        minWidth: 220,
        gap: 4,
    },
    statusPill: {
        borderRadius: 999,
        paddingHorizontal: 10,
        paddingVertical: 6,
    },
    title: {
        fontSize: 18,
        fontWeight: "bold",
    },
    label: {
        fontSize: 14,
    },
    metaGrid: {
        flexDirection: "row",
        gap: 16,
        flexWrap: "wrap",
    },
    metaItem: {
        gap: 4,
        minWidth: 180,
        flex: 1,
    },
    emptyState: {
        borderWidth: 1,
        borderRadius: 18,
        padding: 16,
        gap: 4,
    },
});

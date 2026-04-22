import {View, StyleSheet} from "react-native";
import React, {useState} from "react";
import {LeagueDetailResponse} from "@/hooks/useData";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";
import { AppButton } from "@/components/ui/AppButton";
import AddSeasonModal from "@/components/Seasons/AddSeasonModal";
import usePost from "@/hooks/usePost";

export interface LeagueDetailsCardProps {
    data: LeagueDetailResponse
    onSeasonChanged: () => void
}

const LeagueDetailsCard = (props: LeagueDetailsCardProps) => {
    const theme = useAppTheme();
    const {postData, result, clearResult} = usePost();
    const [openSeasonModal, setOpenSeasonModal] = useState(false);
    const hasActiveSeason = Boolean(props.data.season);
    const seasonStatus = props.data.season ? props.data.season.status : "No active season";
    const canStartSeason = props.data.season?.status === "planned";

    const onCreateSeason = () => {
        setOpenSeasonModal(true);
        clearResult();
    };

    const onSeasonCreated = () => {
        setOpenSeasonModal(false);
        props.onSeasonChanged();
    };

    const onCloseSeasonModal = () => {
        setOpenSeasonModal(false);
    };

    const onStartSeason = async () => {
        if (!props.data.season) {
            return;
        }

        const response = await postData(`/v1/leagues/${props.data.id}/seasons/${props.data.season.id}/start`);
        if (response.ok) {
            props.onSeasonChanged();
        }
    };

    return (
        <View style={styles.content}>
            <View style={styles.header}>
                <View style={styles.headerCopy}>
                    <AppText variant="titleMedium">{props.data.name}</AppText>
                    <AppText variant="bodySmall" style={{color: theme.colors.onSurfaceVariant}}>
                        League details
                    </AppText>
                </View>
                <View
                    style={[
                        styles.statusPill,
                        {
                            backgroundColor: hasActiveSeason ? theme.colors.primaryContainer : theme.colors.surfaceVariant,
                        },
                    ]}
                >
                    <AppText
                        variant="labelMedium"
                        style={{color: hasActiveSeason ? theme.colors.primary : theme.colors.onSurfaceVariant}}
                    >
                        {seasonStatus}
                    </AppText>
                </View>
            </View>

            {result ? (
                <View style={[styles.messagePanel, {backgroundColor: theme.colors.errorContainer, borderColor: theme.colors.error}]}>
                    <AppText variant="bodyMedium" style={{color: theme.colors.onErrorContainer}}>
                        {result}
                    </AppText>
                </View>
            ) : null}

            <View style={styles.section}>
                <AppText variant="labelLarge" style={{color: theme.colors.onSurfaceVariant}}>
                    League
                </AppText>
                <View style={[styles.detailTable, {borderColor: theme.colors.outlineVariant}]}>
                    <DetailRow label="League ID" value={props.data.id} isFirst />
                    <DetailRow label="Active members" value={String(props.data.active_members)} />
                    <DetailRow label="Active season" value={props.data.season ? props.data.season.name : "None"} />
                </View>
            </View>

            <View style={styles.section}>
                <AppText variant="labelLarge" style={{color: theme.colors.onSurfaceVariant}}>
                    Active season
                </AppText>
                {props.data.season ? (
                    <>
                        <View style={[styles.detailTable, {borderColor: theme.colors.outlineVariant}]}>
                            <DetailRow label="Season ID" value={props.data.season.id} isFirst />
                            <DetailRow label="Status" value={props.data.season.status} />
                            <DetailRow label="Phase" value={props.data.season.phase ?? "Regular season"} />
                        </View>
                        {canStartSeason ? (
                            <View style={[styles.actionPanel, {backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant}]}>
                                <View style={styles.actionCopy}>
                                    <AppText variant="titleMedium">Ready to begin</AppText>
                                    <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                                        The schedule is planned. Start the season when fixtures should become active.
                                    </AppText>
                                </View>
                                <AppButton variant="submit" onPress={onStartSeason}>Start Season</AppButton>
                            </View>
                        ) : null}
                    </>
                ) : (
                    <View style={[styles.actionPanel, {backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant}]}>
                        <View style={styles.actionCopy}>
                            <AppText variant="titleMedium">No active season</AppText>
                            <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                                Create a season for this league without leaving the details view.
                            </AppText>
                        </View>
                        <AppButton variant="submit" onPress={onCreateSeason}>Add Season</AppButton>
                    </View>
                )}
            </View>
            <AddSeasonModal onSave={onSeasonCreated} onClose={onCloseSeasonModal} open={openSeasonModal} leagueId={props.data.id}/>
        </View>
    )
}

function DetailRow({label, value, isFirst}: {label: string; value: string; isFirst?: boolean}) {
    const theme = useAppTheme();

    return (
        <View style={[styles.row, {borderTopColor: theme.colors.outlineVariant}, isFirst && styles.firstRow]}>
            <AppText variant="labelMedium" style={[styles.label, {color: theme.colors.onSurfaceVariant}]}>
                {label}
            </AppText>
            <AppText variant="bodyMedium" style={styles.value}>
                {value}
            </AppText>
        </View>
    );
}

export default LeagueDetailsCard;

const styles = StyleSheet.create({
    content: {
        gap: 16,
    },
    header: {
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center",
        gap: 12,
        flexWrap: "wrap",
    },
    headerCopy: {
        gap: 2,
        flex: 1,
        minWidth: 220,
    },
    statusPill: {
        borderRadius: 999,
        paddingHorizontal: 10,
        paddingVertical: 6,
    },
    section: {
        gap: 8,
    },
    messagePanel: {
        borderWidth: 1,
        borderRadius: 16,
        padding: 12,
    },
    detailTable: {
        borderWidth: 1,
        borderRadius: 14,
        overflow: "hidden",
    },
    actionPanel: {
        borderWidth: 1,
        borderRadius: 20,
        padding: 16,
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center",
        flexWrap: "wrap",
        gap: 12,
    },
    actionCopy: {
        flex: 1,
        minWidth: 220,
        gap: 4,
    },
    row: {
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "flex-start",
        gap: 16,
        paddingHorizontal: 14,
        paddingVertical: 12,
        borderTopWidth: 1,
    },
    firstRow: {
        borderTopWidth: 0,
    },
    label: {
        width: 120,
        paddingTop: 2,
    },
    value: {
        flex: 1,
        textAlign: "right",
    }
})

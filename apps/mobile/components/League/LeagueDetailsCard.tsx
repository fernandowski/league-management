import {View, StyleSheet} from "react-native";
import React from "react";
import {LeagueDetailResponse} from "@/hooks/useData";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";

export interface LeagueDetailsCardProps {
    data: LeagueDetailResponse
}

const LeagueDetailsCard = (props: LeagueDetailsCardProps) => {
    const theme = useAppTheme();
    const seasonConfigured = Boolean(props.data.season);
    const seasonStatus = props.data.season ? props.data.season.status : "Not configured";

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
                            backgroundColor: seasonConfigured ? theme.colors.primaryContainer : theme.colors.surfaceVariant,
                        },
                    ]}
                >
                    <AppText
                        variant="labelMedium"
                        style={{color: seasonConfigured ? theme.colors.primary : theme.colors.onSurfaceVariant}}
                    >
                        {seasonStatus}
                    </AppText>
                </View>
            </View>

            <View style={styles.section}>
                <AppText variant="labelLarge" style={{color: theme.colors.onSurfaceVariant}}>
                    League
                </AppText>
                <View style={[styles.detailTable, {borderColor: theme.colors.outlineVariant}]}>
                    <DetailRow label="League ID" value={props.data.id} isFirst />
                    <DetailRow label="Active members" value={String(props.data.active_members)} />
                    <DetailRow label="Season" value={props.data.season ? props.data.season.name : "Not configured"} />
                </View>
            </View>

            <View style={styles.section}>
                <AppText variant="labelLarge" style={{color: theme.colors.onSurfaceVariant}}>
                    Current season
                </AppText>
                {props.data.season ? (
                    <View style={[styles.detailTable, {borderColor: theme.colors.outlineVariant}]}>
                        <DetailRow label="Season ID" value={props.data.season.id} isFirst />
                        <DetailRow label="Status" value={props.data.season.status} />
                    </View>
                ) : (
                    <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                        No season is configured yet. Use the Season tab below to create and plan one.
                    </AppText>
                )}
            </View>
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
    detailTable: {
        borderWidth: 1,
        borderRadius: 14,
        overflow: "hidden",
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

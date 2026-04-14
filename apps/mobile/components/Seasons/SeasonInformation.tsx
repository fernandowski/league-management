import {SeasonDetailResponse} from "@/hooks/useData";
import {StyleSheet, View} from "react-native";
import { AppButton } from "@/components/ui/AppButton";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";

export interface SeasonInformationProps {
    season: SeasonDetailResponse,
    handleSeasonStart: () => void
}
export default function SeasonInformation(props: SeasonInformationProps) {
    const theme = useAppTheme();
    const roundCount = props.season.rounds.length;

    return (
        <View style={styles.surfaceCard}>
            <View style={styles.header}>
                <View style={styles.headerCopy}>
                    <AppText variant="titleLarge">{props.season.name}</AppText>
                    <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                        Track season state, schedule readiness, and round volume from one place.
                    </AppText>
                </View>
                <View style={[styles.statusPill, {backgroundColor: theme.colors.primaryContainer}]}>
                    <AppText variant="labelLarge" style={{color: theme.colors.primary}}>
                        {formatStatus(props.season.status)}
                    </AppText>
                </View>
            </View>

            <View style={styles.metricsRow}>
                <MetricCard label="Status" value={formatStatus(props.season.status)} />
                <MetricCard label="Rounds" value={String(roundCount)} />
            </View>

            {['planned'].indexOf(props.season.status) > -1 && (
                <View style={[styles.actionPanel, {backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant}]}>
                    <View style={styles.actionCopy}>
                        <AppText variant="titleMedium">Ready to begin</AppText>
                        <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                            The schedule is planned. Start the season when fixtures should become active.
                        </AppText>
                    </View>
                    <AppButton variant="submit" onPress={props.handleSeasonStart}>Start Season</AppButton>
                </View>
            )}
        </View>
    )
}

function MetricCard({label, value}: {label: string; value: string}) {
    return (
        <View style={styles.metricCard}>
            <AppText variant="labelLarge" style={styles.metricLabel}>{label}</AppText>
            <AppText variant="headlineSmall">{value}</AppText>
        </View>
    );
}

function formatStatus(status: string) {
    return status.replace(/_/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
}

const styles = StyleSheet.create({
    surfaceCard: {
        gap: 16,
    },
    header: {
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "flex-start",
        gap: 12,
        flexWrap: "wrap",
    },
    headerCopy: {
        flex: 1,
        minWidth: 220,
        gap: 4,
    },
    statusPill: {
        borderRadius: 999,
        paddingHorizontal: 12,
        paddingVertical: 8,
    },
    metricsRow: {
        flexDirection: "row",
        flexWrap: "wrap",
        gap: 12,
    },
    metricCard: {
        flex: 1
    },
    metricLabel: {
        opacity: 0.7,
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
})

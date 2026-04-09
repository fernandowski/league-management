import { useRouter } from "expo-router";
import { useEffect, useMemo, useState } from "react";
import { StyleSheet, View } from "react-native";

import { apiRequest } from "@/api/api";
import ViewContent from "@/components/Layout/ViewContent";
import AddOrganization from "@/components/Overview/AddOrganization";
import { AppButton } from "@/components/ui/AppButton";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import { useOrganizationStore } from "@/stores/organizationStore";
import { useAppTheme } from "@/theme/theme";

interface LeagueResponse {
    data: { id: string; name: string }[];
    total: number;
}

export default function Overview() {
    const router = useRouter();
    const theme = useAppTheme();
    const { organization, organizations } = useOrganizationStore();
    const [leagueCount, setLeagueCount] = useState(0);
    const [teamCount, setTeamCount] = useState(0);
    const [primaryLeagueId, setPrimaryLeagueId] = useState<string | null>(null);

    const selectedOrganization = useMemo(
        () => organizations.find((currentOrganization) => currentOrganization.id === organization) ?? null,
        [organization, organizations]
    );

    useEffect(() => {
        const fetchOverviewData = async () => {
            if (!organization) {
                setLeagueCount(0);
                setTeamCount(0);
                setPrimaryLeagueId(null);
                return;
            }

            try {
                const [leagueResponse, teamResponse] = await Promise.all([
                    apiRequest(
                        `/v1/leagues?offset=0&limit=0&organization_id=${organization}`,
                        { method: "GET" }
                    ) as Promise<LeagueResponse>,
                    apiRequest(`/v1/teams/?organization_id=${organization}`, { method: "GET" }) as Promise<{ id: string }[]>,
                ]);

                setLeagueCount(leagueResponse.total ?? 0);
                setPrimaryLeagueId(leagueResponse.data[0]?.id ?? null);
                setTeamCount(teamResponse.length);
            } catch {
                setLeagueCount(0);
                setTeamCount(0);
                setPrimaryLeagueId(null);
            }
        };

        void fetchOverviewData();
    }, [organization]);

    return (
        <ViewContent>
            <View style={styles.header}>
                <View style={styles.headerTopRow}>
                    <AppText variant={"headlineMedium"} style={styles.headerTitle}>
                        {selectedOrganization?.name ?? "Organization Home"}
                    </AppText>
                </View>
                <AppText variant={"bodyLarge"} style={{ color: theme.colors.onSurfaceVariant }}>
                    Use this space as the control center for the selected organization.
                </AppText>
            </View>

            <View style={styles.summaryGrid}>
                <AppCard style={styles.workspaceCard}>
                    <AppCard.Content style={styles.workspaceContent}>
                        <View style={styles.workspaceHeader}>
                            <AppText variant={"labelLarge"} style={{ color: theme.colors.onSurfaceVariant }}>Current Workspace</AppText>
                            <AddOrganization
                                buttonLabel="Add Organization"
                                description=""
                                style={styles.inlineAddOrganization}
                            />
                        </View>
                        <AppText variant={"headlineSmall"}>{selectedOrganization?.name ?? "No Organization Selected"}</AppText>
                        <AppText variant={"bodyMedium"} style={{ color: theme.colors.onSurfaceVariant }}>
                            This is the organization currently active across leagues, teams, and seasons.
                        </AppText>

                        <View style={styles.workspaceSections}>
                            <AppCard style={[styles.summaryCard, { backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant }]}>
                                <AppCard.Content style={styles.summaryContent}>
                                    <AppText variant={"labelLarge"} style={{ color: theme.colors.onSurfaceVariant }}>Leagues</AppText>
                                    <AppText variant={"headlineMedium"}>{leagueCount}</AppText>
                                    <AppText variant={"bodyMedium"} style={{ color: theme.colors.onSurfaceVariant }}>
                                        In the selected organization
                                    </AppText>
                                    <AppButton
                                        mode={"contained"}
                                        onPress={() => router.push(`/dashboard/leagues`)}
                                    >
                                        {leagueCount > 0 ? "Manage Leagues" : "Create League"}
                                    </AppButton>
                                </AppCard.Content>
                            </AppCard>

                            <AppCard style={[styles.summaryCard, { backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant }]}>
                                <AppCard.Content style={styles.summaryContent}>
                                    <AppText variant={"labelLarge"} style={{ color: theme.colors.onSurfaceVariant }}>Teams</AppText>
                                    <AppText variant={"headlineMedium"}>{teamCount}</AppText>
                                    <AppText variant={"bodyMedium"} style={{ color: theme.colors.onSurfaceVariant }}>
                                        Ready to assign and manage
                                    </AppText>
                                    <AppButton mode={"outlined"} onPress={() => router.push("/dashboard/teams")}>
                                        Manage Teams
                                    </AppButton>
                                </AppCard.Content>
                            </AppCard>

                            <AppCard style={[styles.summaryCard, { backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant }]}>
                                <AppCard.Content style={styles.summaryContent}>
                                    <AppText variant={"labelLarge"} style={{ color: theme.colors.onSurfaceVariant }}>Seasons</AppText>
                                    <AppText variant={"headlineMedium"}>{leagueCount > 0 ? "Active" : "None"}</AppText>
                                    <AppText variant={"bodyMedium"} style={{ color: theme.colors.onSurfaceVariant }}>
                                        Review schedules and season progress
                                    </AppText>
                                    <AppButton mode={"outlined"} onPress={() => router.push("/dashboard/seasons/index")}>
                                        Manage Seasons
                                    </AppButton>
                                </AppCard.Content>
                            </AppCard>
                        </View>
                    </AppCard.Content>
                </AppCard>
            </View>
        </ViewContent>
    )
}

const styles = StyleSheet.create({
    header: {
        gap: 8,
        marginBottom: 24,
    },
    headerTopRow: {
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
        gap: 12,
    },
    headerTitle: {
        flex: 1,
    },
    inlineAddOrganization: {
        marginBottom: 0,
        alignItems: "flex-end",
        flexShrink: 0,
    },
    summaryGrid: {
        marginBottom: 24,
    },
    workspaceCard: {
        borderRadius: 24,
    },
    workspaceContent: {
        gap: 16,
    },
    workspaceHeader: {
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
        gap: 12,
    },
    workspaceSections: {
        gap: 12,
    },
    summaryCard: {
        borderRadius: 24,
    },
    summaryContent: {
        gap: 16,
    },
});

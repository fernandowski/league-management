import {StyleSheet, View} from "react-native";
import {useEffect} from "react";

import {AppCard} from "@/components/ui/AppCard";
import {AppText} from "@/components/ui/AppText";
import {PlayoffRoundRuleResponse, PlayoffRulesResponse, useData} from "@/hooks/useData";
import PlayoffRulesForm from "@/components/Seasons/PlayoffRulesForm";
import {useAppTheme} from "@/theme/theme";

interface SeasonPlayoffsProps {
    seasonId: string
}

export default function SeasonPlayoffs({seasonId}: SeasonPlayoffsProps) {
    const {fetchData, data, error, fetching} = useData<PlayoffRulesResponse>();
    const theme = useAppTheme();

    useEffect(() => {
        if (seasonId) {
            fetchData(`/v1/seasons/${seasonId}/playoffs/rules`);
        }
    }, [fetchData, seasonId]);

    if (fetching && !data) {
        return (
            <AppCard>
                <AppCard.Content>
                    <AppText>Loading playoff configuration</AppText>
                </AppCard.Content>
            </AppCard>
        );
    }

    return (
        <View style={styles.container}>
            <AppCard>
                <AppCard.Content style={styles.summaryContent}>
                    <View style={styles.header}>
                        <AppText variant="titleLarge">Playoffs</AppText>
                        <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                            Configure the postseason now. Qualification preview and bracket generation land in the next slice.
                        </AppText>
                    </View>

                    {error && (
                        <View style={[styles.messageBox, {backgroundColor: theme.colors.errorContainer, borderColor: theme.colors.error}]}>
                            <AppText style={{color: theme.colors.onErrorContainer}}>{error}</AppText>
                        </View>
                    )}

                    {data?.configured && data.rules ? (
                        <View style={[styles.messageBox, {backgroundColor: theme.colors.secondaryContainer, borderColor: theme.colors.outlineVariant}]}>
                            <AppText variant="labelLarge" style={{color: theme.colors.secondary}}>Configured</AppText>
                            <AppText style={{color: theme.colors.onSecondaryContainer}}>
                                Top {data.rules.qualifier_count} qualify. {summarizeRounds(data.rules.rounds)}
                            </AppText>
                        </View>
                    ) : (
                        <View style={[styles.messageBox, {backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant}]}>
                            <AppText style={{color: theme.colors.onSurfaceVariant}}>
                                No playoff rules have been configured for this season yet.
                            </AppText>
                        </View>
                    )}
                </AppCard.Content>
            </AppCard>

            <PlayoffRulesForm seasonId={seasonId} playoffRules={data} onSaved={() => fetchData(`/v1/seasons/${seasonId}/playoffs/rules`)}/>
        </View>
    );
}

function summarizeRounds(rounds: PlayoffRoundRuleResponse[]) {
    return rounds
        .map((round) => `${round.name} ${round.legs === 1 ? "single-leg" : "two-leg"}`)
        .join(", ");
}

const styles = StyleSheet.create({
    container: {
        gap: 16,
    },
    header: {
        gap: 4,
    },
    summaryContent: {
        gap: 16,
    },
    messageBox: {
        borderWidth: 1,
        borderRadius: 16,
        paddingHorizontal: 14,
        paddingVertical: 12,
        gap: 4,
    },
});

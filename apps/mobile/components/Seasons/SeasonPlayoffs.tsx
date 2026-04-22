import {StyleSheet, View} from "react-native";
import {useEffect, useState} from "react";

import {AppButton} from "@/components/ui/AppButton";
import {AppCard} from "@/components/ui/AppCard";
import {AppText} from "@/components/ui/AppText";
import PlayoffBracketView from "@/components/Seasons/PlayoffBracketView";
import PlayoffLegScoreModal from "@/components/Seasons/PlayoffLegScoreModal";
import PlayoffRulesForm from "@/components/Seasons/PlayoffRulesForm";
import { useNotification } from "@/components/ui/AppNotifications";
import {
    PlayoffBracketResponse,
    PlayoffQualificationPreviewResponse,
    PlayoffTieMatchResponse,
    PlayoffRoundRuleResponse,
    PlayoffRulesResponse,
    useData
} from "@/hooks/useData";
import usePost from "@/hooks/usePost";
import {useAppTheme} from "@/theme/theme";

interface SeasonPlayoffsProps {
    seasonId: string
}

export default function SeasonPlayoffs({seasonId}: SeasonPlayoffsProps) {
    const {fetchData, data, error, fetching} = useData<PlayoffRulesResponse>();
    const {fetchData: fetchQualification, data: qualificationPreview} = useData<PlayoffQualificationPreviewResponse>();
    const {fetchData: fetchBracket, data: bracket} = useData<PlayoffBracketResponse>();
    const {postData, result, clearResult} = usePost();
    const { showNotification } = useNotification();
    const theme = useAppTheme();
    const [selectedTieId, setSelectedTieId] = useState<string | null>(null);
    const [selectedMatch, setSelectedMatch] = useState<PlayoffTieMatchResponse | null>(null);

    useEffect(() => {
        if (seasonId) {
            fetchData(`/v1/seasons/${seasonId}/playoffs/rules`);
        }
    }, [fetchData, seasonId]);

    useEffect(() => {
        if (seasonId && data?.configured && data.rules) {
            fetchQualification(`/v1/seasons/${seasonId}/playoffs/qualification-preview`);
            fetchBracket(`/v1/seasons/${seasonId}/playoffs/bracket`);
        }
    }, [data?.configured, data?.rules, fetchBracket, fetchQualification, seasonId]);

    const onGenerateBracket = async () => {
        clearResult();
        const response = await postData(`/v1/seasons/${seasonId}/playoffs/bracket/generate`);
        if (!response.ok) {
            showNotification("Unable to generate the playoff bracket.", "error");
            return;
        }
        await fetchQualification(`/v1/seasons/${seasonId}/playoffs/qualification-preview`);
        await fetchBracket(`/v1/seasons/${seasonId}/playoffs/bracket`);
        showNotification("Playoff bracket generated.", "info");
    };

    const onPlayoffMatchPress = (tieId: string, match: PlayoffTieMatchResponse) => {
        setSelectedTieId(tieId);
        setSelectedMatch(match);
    };

    const onPlayoffMatchSave = async (scores: {home: number; away: number}) => {
        if (!selectedTieId || !selectedMatch) {
            return;
        }

        clearResult();
        const response = await postData(`/v1/seasons/${seasonId}/playoffs/ties/${selectedTieId}/matches/${selectedMatch.id}/score`, {
            match_id: selectedMatch.id,
            home_score: scores.home,
            away_score: scores.away,
        }, "PUT");
        if (!response.ok) {
            if (response.error === "final must have a winner") {
                showNotification("The final cannot end level. Enter a winning score.", "error");
            } else {
                showNotification(response.error || "Unable to save the playoff score.", "error");
            }
            return;
        }
        await fetchBracket(`/v1/seasons/${seasonId}/playoffs/bracket`);
        setSelectedTieId(null);
        setSelectedMatch(null);
        showNotification("Playoff score saved.", "info");
    };

    if (fetching && !data) {
        return (
            <AppCard>
                <AppCard.Content>
                    <AppText>Loading playoff configuration</AppText>
                </AppCard.Content>
            </AppCard>
        );
    }

    const shouldShowBracket = Boolean(bracket?.generated) && qualificationPreview?.bracket_exists !== false;

    return (
        <View style={styles.container}>
            <AppCard>
                <AppCard.Content style={styles.summaryContent}>
                    <View style={styles.header}>
                        <AppText variant="titleLarge">Playoffs</AppText>
                        <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                            Configure the postseason, generate the bracket, and record playoff scores.
                        </AppText>
                    </View>

                    {error && (
                        <View style={[styles.messageBox, {backgroundColor: theme.colors.errorContainer, borderColor: theme.colors.error}]}>
                            <AppText style={{color: theme.colors.onErrorContainer}}>{error}</AppText>
                        </View>
                    )}
                    {result && (
                        <View style={[styles.messageBox, {backgroundColor: theme.colors.errorContainer, borderColor: theme.colors.error}]}>
                            <AppText style={{color: theme.colors.onErrorContainer}}>{result}</AppText>
                        </View>
                    )}

                    {data?.configured && data.rules ? (
                        <View style={[styles.messageBox, {backgroundColor: theme.colors.secondaryContainer, borderColor: theme.colors.outlineVariant}]}>
                            <AppText variant="labelLarge" style={{color: theme.colors.secondary}}>Configured</AppText>
                            <AppText style={{color: theme.colors.onSecondaryContainer}}>
                                {summarizePlayoffRules(data.rules.qualifier_count, data.rules.rounds)}
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

            <PlayoffRulesForm
                seasonId={seasonId}
                playoffRules={data}
                onSaved={async () => {
                    await fetchData(`/v1/seasons/${seasonId}/playoffs/rules`);
                    await fetchQualification(`/v1/seasons/${seasonId}/playoffs/qualification-preview`);
                    await fetchBracket(`/v1/seasons/${seasonId}/playoffs/bracket`);
                    if (data?.bracket_generated && !data?.rules_locked) {
                        showNotification("Playoff rules updated. Regenerate the bracket before playing matches.", "warning");
                    } else {
                        showNotification("Playoff rules saved.", "info");
                    }
                }}
            />

            {qualificationPreview && (
                <AppCard>
                    <AppCard.Content style={styles.summaryContent}>
                        <View style={styles.header}>
                            <AppText variant="titleMedium">Qualification Preview</AppText>
                            <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                                Based on the current final standings, these teams would enter the postseason.
                            </AppText>
                        </View>

                        <View style={styles.qualifierList}>
                            {qualificationPreview.qualified_teams.map((team) => (
                                <View key={team.team_id} style={[styles.qualifierRow, {borderColor: theme.colors.outlineVariant}]}>
                                    <AppText variant="labelLarge">Seed {team.seed}</AppText>
                                    <AppText variant="bodyMedium">{team.team_name}</AppText>
                                    <AppText variant="bodySmall" style={{color: theme.colors.onSurfaceVariant}}>
                                        {team.points} pts
                                    </AppText>
                                </View>
                            ))}
                        </View>

                        {!qualificationPreview.bracket_exists && (
                            <View style={styles.actions}>
                                <AppButton variant="submit" onPress={onGenerateBracket}>
                                    Generate Bracket
                                </AppButton>
                            </View>
                        )}
                    </AppCard.Content>
                </AppCard>
            )}

            {shouldShowBracket && bracket && (
                <PlayoffBracketView bracket={bracket} onMatchPress={onPlayoffMatchPress}/>
            )}

            {selectedMatch && (
                <PlayoffLegScoreModal
                    isOpen={Boolean(selectedMatch)}
                    match={selectedMatch}
                    onClose={() => {
                        setSelectedTieId(null);
                        setSelectedMatch(null);
                    }}
                    onSave={onPlayoffMatchSave}
                />
            )}
        </View>
    );
}

function summarizePlayoffRules(qualifierCount: number, rounds: PlayoffRoundRuleResponse[]) {
    const nonFinalRound = rounds.find((round) => round.name !== "final");
    const nonFinalFormat = nonFinalRound?.legs === 2 ? "two legs" : "one leg";

    if (qualifierCount <= 2) {
        return "Top 2 qualify. The playoffs are a single-match final.";
    }

    return `Top ${qualifierCount} qualify. Earlier knockout rounds use ${nonFinalFormat}; the final is always a single match.`;
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
    qualifierList: {
        gap: 10,
    },
    qualifierRow: {
        borderWidth: 1,
        borderRadius: 14,
        paddingHorizontal: 14,
        paddingVertical: 12,
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
        gap: 12,
    },
    actions: {
        flexDirection: "row",
        justifyContent: "flex-end",
    },
});

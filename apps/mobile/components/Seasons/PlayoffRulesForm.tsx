import {useEffect, useState} from "react";
import {StyleSheet, View} from "react-native";

import {AppButton} from "@/components/ui/AppButton";
import {AppCard} from "@/components/ui/AppCard";
import {AppSelect} from "@/components/ui/AppSelect";
import {AppText} from "@/components/ui/AppText";
import {AppTextField} from "@/components/ui/AppTextField";
import {PlayoffRulesResponse} from "@/hooks/useData";
import usePost from "@/hooks/usePost";
import {useAppTheme} from "@/theme/theme";

interface PlayoffRulesFormProps {
    seasonId: string
    playoffRules: PlayoffRulesResponse | null | undefined
    onSaved: () => void
}

const booleanOptions = [
    {label: "No", value: "false"},
    {label: "Yes", value: "true"},
];

const legsOptions = [
    {label: "Single Leg", value: "1"},
    {label: "Two Legs", value: "2"},
];

export default function PlayoffRulesForm({seasonId, playoffRules, onSaved}: PlayoffRulesFormProps) {
    const {postData, result, clearResult} = usePost();
    const theme = useAppTheme();
    const existingRules = playoffRules?.rules;

    const [qualifierCount, setQualifierCount] = useState("4");
    const [semifinalLegs, setSemifinalLegs] = useState("2");
    const [finalLegs, setFinalLegs] = useState("1");
    const [thirdPlaceMatch, setThirdPlaceMatch] = useState("false");
    const [reseedEachRound, setReseedEachRound] = useState("false");
    const [isSaving, setIsSaving] = useState(false);

    useEffect(() => {
        if (!existingRules) {
            return;
        }

        setQualifierCount(String(existingRules.qualifier_count));
        setThirdPlaceMatch(String(existingRules.third_place_match));
        setReseedEachRound(String(existingRules.reseed_each_round));

        const semifinalRule = existingRules.rounds.find((round) => round.name === "semifinal");
        const finalRule = existingRules.rounds.find((round) => round.name === "final");

        if (semifinalRule) {
            setSemifinalLegs(String(semifinalRule.legs));
        }
        if (finalRule) {
            setFinalLegs(String(finalRule.legs));
        }
    }, [existingRules]);

    const isReadOnly = playoffRules?.season_phase === "playoffs" || playoffRules?.season_phase === "completed" || playoffRules?.season_status === "finished";

    const onSave = async () => {
        clearResult();
        setIsSaving(true);
        await postData(`/v1/seasons/${seasonId}/playoffs/rules`, {
            qualification_type: "top_n",
            qualifier_count: Number(qualifierCount),
            reseed_each_round: reseedEachRound === "true",
            third_place_match: thirdPlaceMatch === "true",
            allow_admin_seed_override: false,
            rounds: [
                {
                    name: "semifinal",
                    legs: Number(semifinalLegs),
                    higher_seed_hosts_second_leg: semifinalLegs === "2",
                    tied_aggregate_resolution: "penalties",
                },
                {
                    name: "final",
                    legs: Number(finalLegs),
                    higher_seed_hosts_second_leg: false,
                    tied_aggregate_resolution: "penalties",
                },
            ],
        }, "PUT");
        setIsSaving(false);
        onSaved();
    };

    return (
        <AppCard>
            <AppCard.Content style={styles.content}>
                <View style={styles.header}>
                    <AppText variant="titleMedium">Playoff Rules</AppText>
                    <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                        Configure the postseason format before the playoff phase begins.
                    </AppText>
                </View>

                {isReadOnly && (
                    <View style={[styles.notice, {backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant}]}>
                        <AppText style={{color: theme.colors.onSurfaceVariant}}>
                            Playoff rules are read-only after the playoff phase starts.
                        </AppText>
                    </View>
                )}

                <AppTextField
                    label="Qualifier Count"
                    keyboardType="numeric"
                    value={qualifierCount}
                    onChangeText={setQualifierCount}
                    disabled={isReadOnly}
                />

                <View style={styles.grid}>
                    <AppSelect
                        label="Semifinal Format"
                        value={semifinalLegs}
                        options={legsOptions}
                        disabled={isReadOnly}
                        onValueChange={setSemifinalLegs}
                    />
                    <AppSelect
                        label="Final Format"
                        value={finalLegs}
                        options={legsOptions}
                        disabled={isReadOnly}
                        onValueChange={setFinalLegs}
                    />
                </View>

                <View style={styles.grid}>
                    <AppSelect
                        label="Third-place Match"
                        value={thirdPlaceMatch}
                        options={booleanOptions}
                        disabled={isReadOnly}
                        onValueChange={setThirdPlaceMatch}
                    />
                    <AppSelect
                        label="Reseed Each Round"
                        value={reseedEachRound}
                        options={booleanOptions}
                        disabled={isReadOnly}
                        onValueChange={setReseedEachRound}
                    />
                </View>

                <View style={[styles.summaryBox, {backgroundColor: theme.colors.primaryContainer}]}>
                    <AppText variant="labelLarge" style={{color: theme.colors.primary}}>MVP Rules Summary</AppText>
                    <AppText style={{color: theme.colors.onPrimaryContainer}}>
                        Qualification uses final standings, ties resolve by penalties, and away goals are not part of this first slice.
                    </AppText>
                </View>

                {result && (
                    <View style={[styles.notice, {backgroundColor: theme.colors.errorContainer, borderColor: theme.colors.error}]}>
                        <AppText style={{color: theme.colors.onErrorContainer}}>{result}</AppText>
                    </View>
                )}

                <View style={styles.actions}>
                    <AppButton variant="submit" onPress={onSave} loading={isSaving} disabled={isSaving || isReadOnly}>
                        Save Playoff Rules
                    </AppButton>
                </View>
            </AppCard.Content>
        </AppCard>
    );
}

const styles = StyleSheet.create({
    content: {
        gap: 16,
    },
    header: {
        gap: 4,
    },
    grid: {
        gap: 12,
    },
    summaryBox: {
        borderRadius: 16,
        padding: 16,
        gap: 8,
    },
    actions: {
        flexDirection: "row",
        justifyContent: "flex-end",
    },
    notice: {
        borderWidth: 1,
        borderRadius: 16,
        paddingHorizontal: 14,
        paddingVertical: 12,
    },
});

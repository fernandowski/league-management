import {StyleSheet, View} from "react-native";
import {useCallback, useEffect, useMemo, useState} from "react";
import {MatchesResponse, MatchScore, useData} from "@/hooks/useData";
import MatchUpRound from "@/components/Seasons/MatchUpRound";
import ManageMatchScoreModal from "@/components/Seasons/ManageMatchScoreModal";
import usePost from "@/hooks/usePost";
import { AppText } from "@/components/ui/AppText";
import MatchesPagination from "@/components/Seasons/MatchesPagination";
import { AppButton } from "@/components/ui/AppButton";
import { useAppTheme } from "@/theme/theme";

export interface SeasonMatchUpManagementProps {
    seasonId: string
    onRoundCompleted?: () => void
    seasonStatus?: string
}

export default function SeasonMatchUpManagement(props: SeasonMatchUpManagementProps) {
    const {fetchData, data, error} = useData<MatchesResponse[]>();
    const [openManageMatch, setOpenMangeMatch] = useState(false);
    const {postData, result, clearResult} = usePost();
    const [matchDetails, setMatchDetails] = useState<MatchScore>();
    const [roundIndex, setRoundIndex] = useState(0);
    const theme = useAppTheme();

    const sortedRounds = useMemo(() => {
        if (!data) {
            return [];
        }

        return [...data].sort((a, b) => a.round - b.round);
    }, [data]);

    const currentRoundIndex = useMemo(() => {
        if (sortedRounds.length === 0) {
            return 0;
        }

        const firstOpenRoundIndex = sortedRounds.findIndex((round) =>
            round.matches.some((match) => match.status !== "finished")
        );

        if (firstOpenRoundIndex > -1) {
            return firstOpenRoundIndex;
        }

        return sortedRounds.length - 1;
    }, [sortedRounds]);

    useEffect(() => {
        setRoundIndex(currentRoundIndex);
    }, [currentRoundIndex]);

    const selectedRound = sortedRounds[roundIndex];
    const isCurrentRound = selectedRound?.round === sortedRounds[currentRoundIndex]?.round;
    const selectedRoundHasOpenMatches = selectedRound?.matches.some((match) => match.status !== "finished") ?? false;
    const seasonAllowsRoundCompletion = props.seasonStatus ? props.seasonStatus === "in_progress" : true;
    const canCompleteSelectedRound = seasonAllowsRoundCompletion && isCurrentRound && selectedRoundHasOpenMatches;

    const onMatchPress = useCallback((matchScoreDetails: MatchScore) => {
        const isByeWeek = matchScoreDetails.home_team === "Bye" || matchScoreDetails.away_team === "Bye";
        if (!canCompleteSelectedRound || isByeWeek || matchScoreDetails.status === "finished") {
            return;
        }
        setMatchDetails(matchScoreDetails);
        setOpenMangeMatch(true);
    }, [canCompleteSelectedRound])

    const onModalClose = () => {
        setOpenMangeMatch(false);
    }

    const onSave = async (scores: {home: number, away: number, matchId: string}) => {
        await postData(`/v1/seasons/${props.seasonId}/matches/score`, {
            match_id: scores.matchId,
            home_score: scores.home,
            away_score: scores.away,
        }, 'PUT');

        await fetchData(`/v1/seasons/${props.seasonId}/matches`);
        onModalClose();
    }

    const onCompleteRound = async () => {
        clearResult();
        await postData(`/v1/seasons/${props.seasonId}/rounds/complete`);
        await fetchData(`/v1/seasons/${props.seasonId}/matches`);
        props.onRoundCompleted?.();
    }

    useEffect(() => {
        if (props.seasonId) {
            fetchData(`/v1/seasons/${props.seasonId}/matches`)
        }

    }, [fetchData, props.seasonId]);

    if (!data) {
        return <View><AppText>{error}</AppText></View>
    }

    if (!selectedRound) {
        return (
            <View>
                <AppText>No matches available.</AppText>
            </View>
        );
    }

    return (
        <View>
            <View style={styles.controlsRow}>
                <MatchesPagination
                    currentRound={selectedRound.round.toString()}
                    onPrevious={() => setRoundIndex((currentIndex) => Math.max(currentIndex - 1, 0))}
                    onNext={() => setRoundIndex((currentIndex) => Math.min(currentIndex + 1, sortedRounds.length - 1))}
                    disablePrevious={roundIndex === 0}
                    disableNext={roundIndex === sortedRounds.length - 1}
                />
                {canCompleteSelectedRound && (
                    <AppButton
                        variant="submit"
                        onPress={onCompleteRound}
                        contentStyle={styles.completeRoundButtonContent}
                        labelStyle={styles.completeRoundButtonLabel}
                    >
                        Complete Round
                    </AppButton>
                )}
            </View>
            {result && (
                <View style={[styles.messageContainer, {backgroundColor: theme.colors.errorContainer}]}>
                    <AppText style={{color: theme.colors.onErrorContainer}}>{result}</AppText>
                </View>
            )}
            <View style={styles.matchUpRoundContainer}>
                <MatchUpRound data={selectedRound} onMatchPress={onMatchPress} isEditableRound={canCompleteSelectedRound}/>
            </View>
            { matchDetails && <ManageMatchScoreModal isOpen={openManageMatch} matchDetails={matchDetails} seasonId={props.seasonId} onSave={onSave} onClose={onModalClose}/>}
        </View>
    )
}

const styles = StyleSheet.create({
    controlsRow: {
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
        gap: 12,
        flexWrap: "wrap",
    },
    messageContainer: {
        marginTop: 12,
        borderRadius: 12,
        paddingHorizontal: 14,
        paddingVertical: 12,
    },
    completeRoundButtonContent: {
        paddingHorizontal: 10,
        minHeight: 44,
    },
    completeRoundButtonLabel: {
        fontSize: 15,
        fontWeight: "700",
        letterSpacing: 0.2,
    },
    matchUpRoundContainer: {
        marginTop: 18
    }
})

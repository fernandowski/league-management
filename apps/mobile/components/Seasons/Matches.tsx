import {SeasonDetailResponse} from "@/hooks/useData";
import {View,StyleSheet} from "react-native";
import MatchesPagination from "@/components/Seasons/MatchesPagination";
import MatchUpList from "@/components/Seasons/MatchUpList";
import {useMemo, useState} from "react";
import { AppText } from "@/components/ui/AppText";
import { AppCard } from "@/components/ui/AppCard";
import { useAppTheme } from "@/theme/theme";

export interface MatchesProps {
    season: SeasonDetailResponse
}
export default function Matches (props: MatchesProps) {
    const [roundNumber, setRoundNumber] = useState(0);
    const theme = useAppTheme();

    const handleHandleNextRound = () => {
        if (props.season.rounds[roundNumber + 1]) {
            setRoundNumber(roundNumber + 1)
        }

    }

    const handlePreviousRound = () => {
        if (props.season.rounds[roundNumber -1]) {
            setRoundNumber(roundNumber - 1);
        }
    }

    const sortedRounds = useMemo(() => {
        return [...props.season.rounds].sort(
            (a, b) => Number(a.round_number) - Number(b.round_number)
        );
    }, [props.season.rounds]);

    if (sortedRounds.length === 0) {
        return (
            <AppCard style={{borderRadius: 24}}>
                <AppCard.Content>
                    <AppText variant="titleMedium">No matches planned yet</AppText>
                </AppCard.Content>
            </AppCard>
        );
    }

    return (
        <AppCard style={[styles.card, {borderColor: theme.colors.outline}]}>
            <AppCard.Content style={styles.content}>
                <View style={styles.header}>
                    <View style={styles.headerCopy}>
                        <AppText variant="titleMedium">Season schedule</AppText>
                        <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                            Review fixtures round by round and move through the schedule quickly.
                        </AppText>
                    </View>
                    <MatchesPagination onNext={handleHandleNextRound} onPrevious={handlePreviousRound} currentRound={sortedRounds[roundNumber].round_number}/>
                </View>
                <View>
                    <MatchUpList round={sortedRounds[roundNumber]}/>
                </View>
            </AppCard.Content>
        </AppCard>
    )
}

const styles = StyleSheet.create({
    card: {
        borderRadius: 24,
    },
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
        flex: 1,
        minWidth: 220,
        gap: 4,
    },
})

import {StyleSheet, View} from "react-native";
import React from "react";
import {LeagueDetailResponse} from "@/hooks/useData";
import LeagueSeasonDetails from "@/components/Seasons/SeasonDetails";
import AddSeason from "@/components/Seasons/AddSeason";
import SeasonsTable from "@/components/Seasons/SeasonsTable";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";

interface SeasonsViewProps {
    league: LeagueDetailResponse
    onSeasonAdded: () =>  void
    onSeasonPlanned: () => void
}

const SeasonTabView = (props: SeasonsViewProps) => {
    const theme = useAppTheme();

    return (
        <View style={styles.container}>
            {props.league.season ? (
                <LeagueSeasonDetails seasonId={props.league.season.id} leagueId={props.league.id} onSeasonPlanned={props.onSeasonPlanned}/>
            ) : (
                <AddSeason leagueId={props.league.id} onSeasonAdded={props.onSeasonAdded}/>
            )}

            <AppCard style={[styles.historyCard, {borderColor: theme.colors.outline}]}>
                <AppCard.Content style={styles.historyContent}>
                    <View style={styles.historyHeader}>
                        <AppText variant="titleLarge">Past seasons</AppText>
                        <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                            Review completed seasons for this league and reopen their details whenever you need historical context.
                        </AppText>
                    </View>
                    <SeasonsTable
                        leagueId={props.league.id}
                        phaseFilter="completed"
                        emptyMessage="Completed seasons for this league will appear here."
                    />
                </AppCard.Content>
            </AppCard>
        </View>
    )
}

export default SeasonTabView

const styles = StyleSheet.create({
    container: {
        flex: 1,
        gap: 16,
    },
    historyCard: {
        borderRadius: 24,
    },
    historyContent: {
        gap: 16,
    },
    historyHeader: {
        gap: 4,
    },
})

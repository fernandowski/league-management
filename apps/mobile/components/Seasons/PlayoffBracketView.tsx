import { StyleSheet, View } from 'react-native';

import { AppButton } from '@/components/ui/AppButton';
import { AppCard } from '@/components/ui/AppCard';
import { AppText } from '@/components/ui/AppText';
import { PlayoffBracketResponse, PlayoffTieMatchResponse } from '@/hooks/useData';
import { useAppTheme } from '@/theme/theme';

interface PlayoffBracketViewProps {
  bracket: PlayoffBracketResponse;
  onMatchPress?: (tieId: string, match: PlayoffTieMatchResponse) => void;
}

export default function PlayoffBracketView({ bracket, onMatchPress }: PlayoffBracketViewProps) {
  const theme = useAppTheme();

  if (!bracket.generated || bracket.rounds.length === 0) {
    return (
      <AppCard>
        <AppCard.Content>
          <AppText variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant }}>
            No playoff bracket has been generated yet.
          </AppText>
        </AppCard.Content>
      </AppCard>
    );
  }

  return (
    <View style={styles.container}>
      {bracket.rounds.map((round) => (
        <AppCard key={`${round.order}-${round.name}`}>
          <AppCard.Content style={styles.roundContent}>
            <View style={styles.roundHeader}>
              <AppText variant="titleMedium">{round.name}</AppText>
              <AppText variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>
                Round {round.order}
              </AppText>
            </View>

            {round.ties.map((tie) => (
              <View key={tie.id} style={[styles.tieCard, { borderColor: theme.colors.outlineVariant }]}>
                <View style={styles.tieHeader}>
                  <AppText variant="labelLarge">Tie {tie.slot_order}</AppText>
                  <AppText variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>
                    {formatTieStatus(tie.status)}
                  </AppText>
                </View>
                <View style={styles.tieRow}>
                  <AppText variant="labelLarge">Seed {tie.home_seed ?? '-'}</AppText>
                  <AppText variant="bodyMedium">{tie.home_team.name}</AppText>
                </View>
                <View style={styles.tieRow}>
                  <AppText variant="labelLarge">Seed {tie.away_seed ?? '-'}</AppText>
                  <AppText variant="bodyMedium">{tie.away_team.name}</AppText>
                </View>

                {tie.winner_team_id && (
                  <View style={[styles.winnerBox, { backgroundColor: theme.colors.primaryContainer }]}>
                    <AppText variant="labelLarge" style={{ color: theme.colors.primary }}>
                      Winner Decided
                    </AppText>
                    <AppText style={{ color: theme.colors.onPrimaryContainer }}>
                      {tie.winner_team_id === tie.home_team.id ? tie.home_team.name : tie.away_team.name}
                    </AppText>
                  </View>
                )}

                <View style={styles.legsSection}>
                  {tie.matches.length === 0 && (
                    <View style={[styles.emptyState, { borderColor: theme.colors.outlineVariant, backgroundColor: theme.colors.surfaceVariant }]}>
                      <AppText variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>
                        {tie.status === 'pending'
                          ? 'This tie is waiting for the previous round to finish before matches can be played.'
                          : 'No playoff matches are available yet for this tie.'}
                      </AppText>
                    </View>
                  )}
                  {tie.matches.map((match) => (
                    <View key={match.id} style={[styles.legRow, { borderColor: theme.colors.outlineVariant }]}>
                      <View style={styles.legCopy}>
                        <AppText variant="labelMedium">Match {match.match_order}</AppText>
                        <AppText variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>
                          {match.home_team} {match.home_score} - {match.away_score} {match.away_team}
                        </AppText>
                      </View>
                      <AppButton
                        variant={match.status === 'finished' ? 'secondary' : 'submit'}
                        disabled={tie.status === 'finished'}
                        onPress={() => onMatchPress?.(tie.id, match)}
                      >
                        {match.status === 'finished' ? 'Edit Score' : 'Enter Score'}
                      </AppButton>
                    </View>
                  ))}
                </View>

              </View>
            ))}
          </AppCard.Content>
        </AppCard>
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    gap: 16,
  },
  roundContent: {
    gap: 12,
  },
  roundHeader: {
    gap: 2,
  },
  tieHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 12,
  },
  tieCard: {
    borderWidth: 1,
    borderRadius: 14,
    padding: 14,
    gap: 10,
  },
  tieRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 12,
  },
  legsSection: {
    gap: 10,
  },
  legRow: {
    borderWidth: 1,
    borderRadius: 12,
    paddingHorizontal: 12,
    paddingVertical: 10,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 12,
  },
  legCopy: {
    flex: 1,
    gap: 2,
  },
  emptyState: {
    borderWidth: 1,
    borderRadius: 12,
    paddingHorizontal: 12,
    paddingVertical: 10,
  },
  winnerBox: {
    borderRadius: 12,
    paddingHorizontal: 12,
    paddingVertical: 10,
    gap: 2,
  },
});

function formatTieStatus(status: string) {
  return status.replace(/_/g, ' ').replace(/\b\w/g, (letter) => letter.toUpperCase());
}

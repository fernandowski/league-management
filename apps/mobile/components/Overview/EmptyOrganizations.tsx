import { StyleSheet, View } from 'react-native';

import AddOrganization from '@/components/Overview/AddOrganization';
import { AppCard } from '@/components/ui/AppCard';
import { AppScreen } from '@/components/ui/AppScreen';
import { AppText } from '@/components/ui/AppText';
import { themeTokens, useAppTheme } from '@/theme/theme';

export function EmptyOrganizations() {
  const theme = useAppTheme();

  return (
    <AppScreen style={styles.screen}>
      <View style={styles.topOrb} />
      <View style={styles.bottomOrb} />

      <AppCard style={styles.card}>
        <AppCard.Content style={styles.content}>
          <View style={styles.copy}>
            <AppText variant="headlineMedium">Start with your first organization</AppText>
            <AppText variant="bodyLarge" style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
              This should be a dedicated first-run step. Until an organization exists, the dashboard navigation only
              adds noise.
            </AppText>
            <AppText variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant }}>
              Create the organization that owns your leagues, teams, and seasons, then the full workspace can open.
            </AppText>
          </View>

          <AddOrganization
            buttonLabel="Create Organization"
            description="Give your workspace a name to unlock leagues, teams, and seasons."
            modalTitle="Create your first organization"
            style={styles.action}
          />
        </AppCard.Content>
      </AppCard>
    </AppScreen>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  topOrb: {
    position: 'absolute',
    top: 88,
    right: 24,
    width: 180,
    height: 180,
    borderRadius: 90,
    backgroundColor: themeTokens.colors.secondarySoft,
  },
  bottomOrb: {
    position: 'absolute',
    bottom: 72,
    left: 12,
    width: 140,
    height: 140,
    borderRadius: 70,
    backgroundColor: themeTokens.colors.primarySoft,
  },
  card: {
    width: '100%',
    maxWidth: 640,
    borderRadius: 32,
  },
  content: {
    gap: 28,
    paddingVertical: 20,
  },
  copy: {
    gap: 12,
  },
  subtitle: {
    maxWidth: 520,
  },
  action: {
    gap: 12,
    marginBottom: 0,
  },
});

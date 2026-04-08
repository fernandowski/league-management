import Ionicons from '@expo/vector-icons/Ionicons';
import { useNavigation } from 'expo-router';
import { useEffect } from 'react';
import { StyleSheet, useWindowDimensions, View } from 'react-native';
import type { DrawerNavigationProp } from '@react-navigation/drawer';
import { IconButton } from 'react-native-paper';

import { Select } from '@/components/Select/Select';
import { AppText } from '@/components/ui/AppText';
import { useOrganizationStore } from '@/stores/organizationStore';
import { useAppTheme } from '@/theme/theme';

export function Header() {
  const navigation: DrawerNavigationProp<any> = useNavigation();
  const dimensions = useWindowDimensions();
  const isLargeScreen = dimensions.width >= 768;
  const theme = useAppTheme();

  const { organizations, fetchOrganizations, setOrganization, organization } = useOrganizationStore();

  useEffect(() => {
    fetchOrganizations();
  }, [fetchOrganizations]);

  return (
    <View
      style={[
        styles.container,
        {
          backgroundColor: theme.colors.surface,
          borderColor: theme.colors.outline,
          shadowColor: theme.colors.shadow,
        },
      ]}
    >
      {!isLargeScreen && (
        <IconButton
          accessibilityLabel="Open navigation menu"
          icon={() => <Ionicons name="menu" size={20} color={theme.colors.primary} />}
          onPress={() => navigation.openDrawer()}
          mode="contained-tonal"
          containerColor={theme.colors.primaryContainer}
          style={styles.menuButton}
        />
      )}

      <View style={styles.content}>
        <AppText variant="labelLarge" style={[styles.label, { color: theme.colors.onSurfaceVariant }]}>
          Organization
        </AppText>
        <Select
          onChange={setOrganization}
          selected={organization}
          data={organizations.map((currentOrganization) => ({
            label: currentOrganization.name,
            value: currentOrganization.id,
          }))}
        />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    minHeight: 84,
    borderWidth: 1,
    borderRadius: 24,
    paddingHorizontal: 16,
    paddingVertical: 14,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.08,
    shadowRadius: 16,
    elevation: 4,
  },
  menuButton: {
    margin: 0,
  },
  content: {
    flex: 1,
    gap: 8,
  },
  label: {
    fontSize: 13,
    letterSpacing: 1.2,
    textTransform: 'uppercase',
    fontWeight: '600',
  },
});

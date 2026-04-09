import Ionicons from '@expo/vector-icons/Ionicons';
import { useNavigation } from 'expo-router';
import { useEffect } from 'react';
import { StyleSheet, View } from 'react-native';
import type { DrawerNavigationProp } from '@react-navigation/drawer';
import { IconButton } from 'react-native-paper';

import { Select } from '@/components/Select/Select';
import { AppText } from '@/components/ui/AppText';
import { useOrganizationStore } from '@/stores/organizationStore';
import { useAppTheme } from '@/theme/theme';

export function Header() {
  const navigation: DrawerNavigationProp<any> = useNavigation();
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
        },
      ]}
    >
      <IconButton
        accessibilityLabel="Toggle navigation menu"
        icon={() => <Ionicons name="menu" size={20} color={theme.colors.primary} />}
        onPress={() => navigation.toggleDrawer()}
        mode="contained-tonal"
        containerColor={theme.colors.primaryContainer}
        style={styles.menuButton}
      />

      <View style={styles.row}>
        <AppText variant="titleMedium" style={[styles.title, { color: theme.colors.onSurface }]}>
          Organization
        </AppText>

        <View style={styles.selectWrap}>
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
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    minHeight: 72,
    paddingVertical: 14,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  menuButton: {
    margin: 0,
  },
  row: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  title: {
    flexShrink: 0,
    fontSize: 18,
  },
  selectWrap: {
    flex: 1,
  },
});

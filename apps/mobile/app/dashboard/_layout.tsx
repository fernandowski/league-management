import Ionicons from '@expo/vector-icons/Ionicons';
import { Redirect, Slot, usePathname } from 'expo-router';
import { Drawer } from 'expo-router/drawer';
import { useEffect } from 'react';
import { ActivityIndicator } from 'react-native-paper';
import { useWindowDimensions } from 'react-native';
import { View } from 'react-native';
import { GestureHandlerRootView } from 'react-native-gesture-handler';

import { useOrganizationStore } from '@/stores/organizationStore';
import { useAppTheme } from '@/theme/theme';

export default function Layout() {
  const dimensions = useWindowDimensions();
  const isLargeScreen = dimensions.width >= 768;
  const theme = useAppTheme();
  const pathname = usePathname();
  const { fetchOrganizations, loading, initialized, organizations } = useOrganizationStore();

  useEffect(() => {
    void fetchOrganizations();
  }, [fetchOrganizations]);

  if (!initialized) {
    return (
      <View
        style={{
          flex: 1,
          alignItems: 'center',
          justifyContent: 'center',
          backgroundColor: theme.colors.background,
        }}
      >
        <ActivityIndicator size="large" color={theme.colors.primary} />
      </View>
    );
  }

  if (organizations.length === 0) {
    if (pathname !== '/dashboard') {
      return <Redirect href="/dashboard" />;
    }

    return <Slot />;
  }

  return (
    <GestureHandlerRootView style={{ flex: 1 }}>
      <Drawer
        screenOptions={{
          drawerType: isLargeScreen ? 'front' : undefined,
          header: () => null,
          drawerPosition: 'left',
          unmountOnBlur: true,
          sceneContainerStyle: {
            backgroundColor: theme.colors.background,
          },
          drawerStyle: {
            backgroundColor: theme.colors.surface,
            borderRightWidth: 1,
            borderRightColor: theme.colors.outline,
            width: isLargeScreen ? 280 : undefined,
          },
          drawerItemStyle: {
            borderRadius: 14,
            marginHorizontal: 12,
          },
          drawerActiveTintColor: theme.colors.primary,
          drawerInactiveTintColor: theme.colors.onSurfaceVariant,
          drawerActiveBackgroundColor: theme.colors.primaryContainer,
          drawerLabelStyle: {
            marginLeft: 0,
          },
          drawerIcon: ({ color, size }) => <Ionicons name="ellipse-outline" size={size} color={color} />,
        }}
      >
        <Drawer.Screen
          name="index"
          options={{
            drawerLabel: 'Organizations',
            title: 'Organizations',
            drawerIcon: ({ color, size }) => <Ionicons name="business-outline" size={size} color={color} />,
          }}
        />
        <Drawer.Screen
          name="leagues/index"
          options={{
            drawerLabel: 'Leagues',
            title: 'Leagues',
            drawerIcon: ({ color, size }) => <Ionicons name="trophy-outline" size={size} color={color} />,
          }}
        />
        <Drawer.Screen
          name="leagues/[id]/index"
          options={{
            drawerItemStyle: { display: 'none' },
            drawerLabel: () => null,
          }}
        />
        <Drawer.Screen
          name="teams/index"
          options={{
            drawerLabel: 'Teams',
            title: 'Teams',
            drawerIcon: ({ color, size }) => <Ionicons name="people-outline" size={size} color={color} />,
          }}
        />
        <Drawer.Screen
          name="seasons/index"
          options={{
            drawerLabel: 'Seasons',
            title: 'Seasons',
            drawerIcon: ({ color, size }) => <Ionicons name="calendar-outline" size={size} color={color} />,
          }}
        />
        <Drawer.Screen
          name="seasons/[id]"
          options={{
            drawerItemStyle: { display: 'none' },
            drawerLabel: () => null,
          }}
        />
      </Drawer>
    </GestureHandlerRootView>
  );
}

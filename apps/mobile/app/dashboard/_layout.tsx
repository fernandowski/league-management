import { Drawer } from 'expo-router/drawer';
import { useWindowDimensions } from 'react-native';
import { GestureHandlerRootView } from 'react-native-gesture-handler';

import { useAppTheme } from '@/theme/theme';

export default function Layout() {
  const dimensions = useWindowDimensions();
  const isLargeScreen = dimensions.width >= 768;
  const theme = useAppTheme();

  return (
    <GestureHandlerRootView style={{ flex: 1 }}>
      <Drawer
        screenOptions={{
          drawerType: isLargeScreen ? 'permanent' : undefined,
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
          },
          drawerActiveTintColor: theme.colors.primary,
          drawerInactiveTintColor: theme.colors.onSurfaceVariant,
          drawerActiveBackgroundColor: theme.colors.primaryContainer,
          drawerLabelStyle: {
            marginLeft: -12,
          },
        }}
      >
        <Drawer.Screen
          name="index"
          options={{
            drawerLabel: 'Organizations Overview',
            title: 'Organizations Overview',
          }}
        />
        <Drawer.Screen
          name="leagues/[id]/index"
          options={{
            drawerLabel: 'Leagues',
            title: 'Leagues',
          }}
        />
        <Drawer.Screen
          name="teams/index"
          options={{
            drawerLabel: 'Teams',
            title: 'Teams',
          }}
        />
        <Drawer.Screen
          name="seasons/index"
          options={{
            drawerLabel: 'Seasons',
            title: 'Seasons',
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

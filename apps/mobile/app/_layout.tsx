import { ThemeProvider } from '@react-navigation/native';
import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { PaperProvider } from 'react-native-paper';

import { AppNotificationsProvider } from '@/components/ui/AppNotifications';
import { navigationTheme, theme } from '@/theme/theme';

export default function RootLayout() {
  return (
    <PaperProvider theme={theme}>
      <AppNotificationsProvider>
        <ThemeProvider value={navigationTheme}>
          <StatusBar style="dark" />
          <Stack
            screenOptions={{
              headerShown: false,
              contentStyle: {
                backgroundColor: theme.colors.background,
              },
            }}
          />
        </ThemeProvider>
      </AppNotificationsProvider>
    </PaperProvider>
  );
}

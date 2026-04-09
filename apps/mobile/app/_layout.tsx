import { ThemeProvider } from '@react-navigation/native';
import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { PaperProvider } from 'react-native-paper';

import { navigationTheme, theme } from '@/theme/theme';

export default function RootLayout() {
  return (
    <PaperProvider theme={theme}>
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
    </PaperProvider>
  );
}

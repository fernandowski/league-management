import { DefaultTheme as NavigationDefaultTheme, Theme as NavigationTheme } from '@react-navigation/native';
import { Platform } from 'react-native';
import {
  adaptNavigationTheme,
  configureFonts,
  MD3LightTheme,
  type MD3Theme,
  useTheme,
} from 'react-native-paper';

const palette = {
  ink: '#14281f',
  slate: '#5d7167',
  mist: '#eef6f1',
  paper: '#fbfdfb',
  line: '#d3e2d8',
  brand: '#2f7d4a',
  brandDeep: '#1f5a34',
  brandSoft: '#dcefe2',
  accent: '#8bcf9f',
  accentSoft: '#edf8f0',
  success: '#2f7d4a',
  danger: '#c43d4b',
  shadow: '#102019',
} as const;

export const themeTokens = {
  colors: {
    background: palette.mist,
    backgroundAccent: palette.paper,
    surface: '#ffffff',
    surfaceMuted: '#f4f7fb',
    surfaceElevated: '#ffffff',
    surfaceInverse: palette.ink,
    text: palette.ink,
    textMuted: palette.slate,
    textInverse: '#ffffff',
    primary: palette.brand,
    primaryStrong: palette.brandDeep,
    primarySoft: palette.brandSoft,
    secondary: palette.accent,
    secondarySoft: palette.accentSoft,
    border: palette.line,
    borderStrong: '#afc0d1',
    success: palette.success,
    error: palette.danger,
    shadow: palette.shadow,
  },
  spacing: {
    xxs: 4,
    xs: 8,
    sm: 12,
    md: 16,
    lg: 24,
    xl: 32,
    xxl: 48,
  },
  radius: {
    sm: 10,
    md: 16,
    lg: 18,
    pill: 999,
  },
  elevation: {
    card: 4,
    modal: 10,
  },
  typography: {
    display: Platform.select({
      ios: 'System',
      android: 'sans-serif-medium',
      default: 'System',
    }),
    body: Platform.select({
      ios: 'System',
      android: 'sans-serif',
      default: 'System',
    }),
  },
} as const;

const fontConfig = configureFonts({
  config: {
    displayLarge: {
      fontFamily: themeTokens.typography.display,
      letterSpacing: -0.4,
      fontWeight: '700',
    },
    displayMedium: {
      fontFamily: themeTokens.typography.display,
      letterSpacing: -0.3,
      fontWeight: '700',
    },
    displaySmall: {
      fontFamily: themeTokens.typography.display,
      letterSpacing: -0.2,
      fontWeight: '700',
    },
    headlineLarge: {
      fontFamily: themeTokens.typography.display,
      letterSpacing: -0.2,
      fontWeight: '700',
    },
    headlineMedium: {
      fontFamily: themeTokens.typography.display,
      letterSpacing: -0.15,
      fontWeight: '700',
    },
    headlineSmall: {
      fontFamily: themeTokens.typography.display,
      letterSpacing: -0.1,
      fontWeight: '700',
    },
    titleLarge: {
      fontFamily: themeTokens.typography.display,
      letterSpacing: -0.05,
      fontWeight: '700',
    },
    titleMedium: {
      fontFamily: themeTokens.typography.body,
      fontWeight: '600',
    },
    titleSmall: {
      fontFamily: themeTokens.typography.body,
      fontWeight: '600',
    },
    labelLarge: {
      fontFamily: themeTokens.typography.body,
      fontWeight: '600',
    },
    labelMedium: {
      fontFamily: themeTokens.typography.body,
      fontWeight: '500',
    },
    labelSmall: {
      fontFamily: themeTokens.typography.body,
      fontWeight: '500',
    },
    bodyLarge: {
      fontFamily: themeTokens.typography.body,
      fontWeight: '400',
    },
    bodyMedium: {
      fontFamily: themeTokens.typography.body,
      fontWeight: '400',
    },
    bodySmall: {
      fontFamily: themeTokens.typography.body,
      fontWeight: '400',
    },
  },
});

const baseTheme = {
  ...MD3LightTheme,
  fonts: fontConfig,
  roundness: themeTokens.radius.md,
  colors: {
    ...MD3LightTheme.colors,
    primary: themeTokens.colors.primary,
    onPrimary: themeTokens.colors.textInverse,
    primaryContainer: themeTokens.colors.primarySoft,
    onPrimaryContainer: themeTokens.colors.primaryStrong,
    secondary: themeTokens.colors.secondary,
    onSecondary: themeTokens.colors.text,
    secondaryContainer: themeTokens.colors.secondarySoft,
    onSecondaryContainer: themeTokens.colors.text,
    tertiary: themeTokens.colors.primaryStrong,
    onTertiary: themeTokens.colors.textInverse,
    tertiaryContainer: themeTokens.colors.primarySoft,
    onTertiaryContainer: themeTokens.colors.primaryStrong,
    error: themeTokens.colors.error,
    onError: themeTokens.colors.textInverse,
    errorContainer: '#ffd9dd',
    onErrorContainer: '#7d1a28',
    background: themeTokens.colors.background,
    onBackground: themeTokens.colors.text,
    surface: themeTokens.colors.surface,
    onSurface: themeTokens.colors.text,
    surfaceVariant: themeTokens.colors.surfaceMuted,
    onSurfaceVariant: themeTokens.colors.textMuted,
    outline: themeTokens.colors.border,
    outlineVariant: themeTokens.colors.borderStrong,
    inverseSurface: themeTokens.colors.surfaceInverse,
    inverseOnSurface: themeTokens.colors.textInverse,
    inversePrimary: '#9fe2b2',
    shadow: themeTokens.colors.shadow,
    scrim: themeTokens.colors.shadow,
    backdrop: 'rgba(16, 32, 51, 0.28)',
  },
  custom: themeTokens,
} satisfies MD3Theme & { custom: typeof themeTokens };

const navigationThemes = adaptNavigationTheme({
  reactNavigationLight: NavigationDefaultTheme,
  materialLight: baseTheme,
});

export const theme = baseTheme;
export type AppTheme = typeof theme;

export const navigationTheme: NavigationTheme = {
  ...navigationThemes.LightTheme,
  colors: {
    ...navigationThemes.LightTheme.colors,
    background: theme.colors.background,
    card: theme.colors.surface,
    text: theme.colors.onSurface,
    border: theme.colors.outline,
    primary: theme.colors.primary,
    notification: theme.colors.secondary,
  },
};

export function useAppTheme() {
  return useTheme<AppTheme>();
}

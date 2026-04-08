import { PropsWithChildren } from 'react';
import { StyleProp, StyleSheet, View, ViewStyle } from 'react-native';

import { useAppTheme } from '@/theme/theme';

type AppScreenProps = PropsWithChildren<{
  style?: StyleProp<ViewStyle>;
}>;

export function AppScreen({ children, style }: AppScreenProps) {
  const theme = useAppTheme();

  return <View style={[styles.screen, { backgroundColor: theme.colors.background }, style]}>{children}</View>;
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
  },
});

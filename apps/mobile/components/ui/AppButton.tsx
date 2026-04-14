import { ComponentProps } from 'react';
import { StyleSheet } from 'react-native';
import { Button } from 'react-native-paper';

import { useAppTheme } from '@/theme/theme';

export type AppButtonVariant = 'submit' | 'destructive' | 'secondary' | 'ghost' | 'default';

export interface AppButtonProps extends ComponentProps<typeof Button> {
  variant?: AppButtonVariant;
}

export function AppButton({ variant, mode, buttonColor, textColor, style, contentStyle, ...props }: AppButtonProps) {
  const theme = useAppTheme();
  const semanticVariant = getVariantStyles(variant, theme);

  return (
    <Button
      mode={mode ?? semanticVariant.mode}
      buttonColor={buttonColor ?? semanticVariant.buttonColor}
      textColor={textColor ?? semanticVariant.textColor}
      style={[styles.button, semanticVariant.style, style]}
      contentStyle={[styles.content, contentStyle]}
      {...props}
    />
  );
}

function getVariantStyles(variant: AppButtonVariant | undefined, theme: ReturnType<typeof useAppTheme>) {
  switch (variant) {
    case 'submit':
      return {
        mode: 'contained' as const,
        buttonColor: theme.colors.primary,
        textColor: theme.colors.onPrimary,
        style: styles.emphasisButton,
      };
    case 'destructive':
      return {
        mode: 'contained' as const,
        buttonColor: theme.colors.error,
        textColor: theme.colors.onError,
        style: styles.emphasisButton,
      };
    case 'secondary':
      return {
        mode: 'outlined' as const,
        buttonColor: undefined,
        textColor: theme.colors.onSurface,
        style: undefined,
      };
    case 'ghost':
      return {
        mode: 'text' as const,
        buttonColor: undefined,
        textColor: theme.colors.onSurfaceVariant,
        style: undefined,
      };
    case 'default':
      return {
        mode: 'contained-tonal' as const,
        buttonColor: theme.colors.surfaceVariant,
        textColor: theme.colors.onSurface,
        style: undefined,
      };
    default:
      return {
        mode: undefined,
        buttonColor: undefined,
        textColor: undefined,
        style: undefined,
      };
  }
}

const styles = StyleSheet.create({
  button: {
    borderRadius: 12,
  },
  content: {
    minHeight: 42,
    paddingHorizontal: 8,
  },
  emphasisButton: {
    shadowOpacity: 0,
  },
});

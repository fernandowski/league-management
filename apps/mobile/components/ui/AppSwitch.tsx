import { ComponentProps } from 'react';
import { Pressable, StyleSheet, View } from 'react-native';
import { Switch } from 'react-native-paper';

import { AppText } from '@/components/ui/AppText';
import { useAppTheme } from '@/theme/theme';

export interface AppSwitchProps extends Omit<ComponentProps<typeof Switch>, 'value' | 'onValueChange'> {
  value: boolean;
  label: string;
  description?: string;
  onValueChange: (value: boolean) => void;
}

export function AppSwitch({
  value,
  label,
  description,
  disabled,
  onValueChange,
  ...props
}: AppSwitchProps) {
  const theme = useAppTheme();

  return (
    <Pressable
      accessibilityRole="switch"
      accessibilityState={{ checked: value, disabled }}
      disabled={disabled}
      onPress={() => {
        if (!disabled) {
          onValueChange(!value);
        }
      }}
      style={[
        styles.container,
        {
          backgroundColor: disabled ? theme.colors.surfaceVariant : theme.colors.surface,
          borderColor: value ? theme.colors.primary : theme.colors.outline,
        },
      ]}
    >
      <View style={styles.copy}>
        <AppText variant="titleSmall">{label}</AppText>
        {description ? (
          <AppText variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>
            {description}
          </AppText>
        ) : null}
      </View>
      <Switch value={value} disabled={disabled} onValueChange={onValueChange} {...props} />
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: {
    minHeight: 72,
    borderWidth: 1,
    borderRadius: 18,
    paddingHorizontal: 16,
    paddingVertical: 14,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 16,
  },
  copy: {
    flex: 1,
    gap: 4,
  },
});

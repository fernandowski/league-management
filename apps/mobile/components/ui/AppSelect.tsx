import Ionicons from '@expo/vector-icons/Ionicons';
import { useMemo, useState } from 'react';
import { Pressable, StyleSheet, View } from 'react-native';
import { Menu } from 'react-native-paper';

import { AppText } from '@/components/ui/AppText';
import { useAppTheme } from '@/theme/theme';

type SelectOption = {
  value: string;
  label: string;
};

interface AppSelectProps {
  options: SelectOption[];
  value?: string | null;
  label?: string;
  placeholder?: string;
  maxWidth?: number;
  onValueChange: (value: string) => void;
}

export function AppSelect({
  options,
  value,
  label,
  placeholder,
  maxWidth = 320,
  onValueChange,
}: AppSelectProps) {
  const [visible, setVisible] = useState(false);
  const theme = useAppTheme();

  const selectedLabel = useMemo(() => {
    return options.find((option) => option.value === value)?.label ?? placeholder ?? 'Select an option';
  }, [options, placeholder, value]);

  const closeMenu = () => {
    setVisible(false);
  };

  return (
    <View style={styles.container}>
      {label ? (
        <AppText variant="labelLarge" style={{ color: theme.colors.onSurfaceVariant }}>
          {label}
        </AppText>
      ) : null}
      <Menu
        visible={visible}
        onDismiss={closeMenu}
        anchorPosition="bottom"
        contentStyle={[
          styles.menuContent,
          {
            backgroundColor: theme.colors.surface,
            borderColor: theme.colors.outline,
          },
        ]}
        anchor={
          <Pressable
            accessibilityRole="button"
            onPress={() => setVisible((currentValue) => !currentValue)}
            style={({ pressed }) => [
              styles.shell,
              {
                backgroundColor: theme.colors.surface,
                borderColor: visible ? theme.colors.primary : theme.colors.outline,
                shadowColor: theme.colors.shadow,
                maxWidth,
              },
              pressed && { backgroundColor: theme.colors.surfaceVariant },
            ]}
          >
            <AppText numberOfLines={1} style={[styles.value, { color: theme.colors.onSurface }]}>
              {selectedLabel}
            </AppText>
            <View style={[styles.chevronWrap, { backgroundColor: theme.colors.primaryContainer }]}>
              <Ionicons name={visible ? 'chevron-up' : 'chevron-down'} size={18} color={theme.colors.primary} />
            </View>
          </Pressable>
        }
      >
        {placeholder ? (
          <Menu.Item
            title={placeholder}
            onPress={() => {
              onValueChange('');
              closeMenu();
            }}
          />
        ) : null}
        {options.map((option) => (
          <Menu.Item
            key={option.value}
            title={option.label}
            onPress={() => {
              onValueChange(option.value);
              closeMenu();
            }}
            trailingIcon={option.value === value ? 'check' : undefined}
          />
        ))}
      </Menu>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    gap: 8,
  },
  shell: {
    borderWidth: 1,
    borderRadius: 18,
    minHeight: 56,
    justifyContent: 'space-between',
    alignItems: 'center',
    flexDirection: 'row',
    paddingLeft: 16,
    paddingRight: 8,
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.08,
    shadowRadius: 14,
    elevation: 3,
  },
  value: {
    flex: 1,
    paddingRight: 12,
    fontSize: 15,
    fontWeight: '600',
  },
  chevronWrap: {
    width: 32,
    height: 32,
    borderRadius: 16,
    alignItems: 'center',
    justifyContent: 'center',
  },
  menuContent: {
    borderWidth: 1,
    borderRadius: 18,
    minWidth: 220,
  },
});

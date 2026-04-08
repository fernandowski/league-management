import { Picker } from '@react-native-picker/picker';
import { useEffect, useState } from 'react';
import { StyleSheet, View } from 'react-native';

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
  onValueChange: (value: string) => void;
}

export function AppSelect({
  options,
  value,
  label,
  placeholder,
  onValueChange,
}: AppSelectProps) {
  const [selectedValue, setSelectedValue] = useState(value ?? '');
  const theme = useAppTheme();

  useEffect(() => {
    setSelectedValue(value ?? '');
  }, [value]);

  return (
    <View style={styles.container}>
      {label ? (
        <AppText variant="labelLarge" style={{ color: theme.colors.onSurfaceVariant }}>
          {label}
        </AppText>
      ) : null}
      <View
        style={[
          styles.shell,
          {
            backgroundColor: theme.colors.surfaceVariant,
            borderColor: theme.colors.outline,
          },
        ]}
      >
        <Picker
          style={[styles.picker, { color: theme.colors.onSurface }]}
          selectedValue={selectedValue}
          onValueChange={(nextValue: string) => {
            setSelectedValue(nextValue);
            onValueChange(nextValue);
          }}
        >
          {placeholder ? <Picker.Item label={placeholder} value="" /> : null}
          {options.map((item) => (
            <Picker.Item label={item.label} value={item.value} key={item.value} />
          ))}
        </Picker>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    gap: 8,
  },
  shell: {
    maxWidth: 280,
    borderWidth: 1,
    borderRadius: 16,
    overflow: 'hidden',
  },
  picker: {
    maxWidth: 280,
  },
});

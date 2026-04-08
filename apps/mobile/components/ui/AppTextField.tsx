import { ComponentProps } from 'react';
import { TextInput } from 'react-native-paper';

import { useAppTheme } from '@/theme/theme';

export type AppTextFieldProps = ComponentProps<typeof TextInput>;

export function AppTextField({ outlineStyle, style, ...props }: AppTextFieldProps) {
  const theme = useAppTheme();

  return (
    <TextInput
      mode="outlined"
      outlineStyle={[
        {
          borderRadius: 16,
          borderColor: theme.colors.outline,
        },
        outlineStyle,
      ]}
      style={[
        {
          backgroundColor: theme.colors.surface,
        },
        style,
      ]}
      {...props}
    />
  );
}

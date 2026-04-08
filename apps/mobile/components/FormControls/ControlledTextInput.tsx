import { StyleProp, View, ViewStyle } from 'react-native';
import { Control, Controller, FieldValues, Path } from 'react-hook-form';
import { HelperText } from 'react-native-paper';

import { AppTextField } from '@/components/ui/AppTextField';

interface FieldProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  value?: string;
  label?: string;
  rules?: object;
  style?: StyleProp<ViewStyle>;
  [key: string]: any;
}

export default function ControlledTextInput<T extends FieldValues>({
  style,
  control,
  name,
  value = '',
  label,
  rules = {},
  ...rest
}: FieldProps<T>) {
  return (
    <Controller
      control={control}
      name={name}
      defaultValue={value as any}
      rules={rules}
      render={({ field: { onChange, onBlur, value: fieldValue }, fieldState: { error } }) => (
        <View>
          <AppTextField
            label={label}
            value={fieldValue ?? ''}
            onBlur={onBlur}
            onChangeText={onChange}
            style={style}
            {...rest}
          />
          <HelperText type="error" visible={Boolean(error)}>
            {error?.message}
          </HelperText>
        </View>
      )}
    />
  );
}

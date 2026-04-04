import {Control, Controller, FieldValues, Path} from "react-hook-form";
import {TextInput, Text} from "react-native-paper";
import {StyleProp, ViewStyle} from "react-native";

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
                                       style, control, name, value = '', label, rules = {}, ...rest
                                   }: FieldProps<T>) {
    return (
        <Controller
            control={control}
            name={name}
            rules={rules}
            render={({field: {onChange, onBlur, value}, fieldState: { error}}) => (
                <>
                    <TextInput
                        label={label}
                        value={value ?? ''}
                        onBlur={onBlur}
                        onChangeText={onChange}
                        mode={"outlined"}
                        style={style}
                        {...rest}
                    />
                    {error && <Text style={{ color: 'red' }}>{error.message}</Text>}
                </>
            )}>
        </Controller>
    )
}

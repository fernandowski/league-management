import { ComponentProps } from 'react';
import { Button } from 'react-native-paper';

export type AppButtonProps = ComponentProps<typeof Button>;

export function AppButton(props: AppButtonProps) {
  return <Button {...props} />;
}

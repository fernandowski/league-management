import { ComponentProps } from 'react';
import { StyleSheet } from 'react-native';
import { Card } from 'react-native-paper';

import { useAppTheme } from '@/theme/theme';

export type AppCardProps = ComponentProps<typeof Card>;

function AppCardBase({ style, ...props }: AppCardProps) {
  const theme = useAppTheme();

  return (
    <Card
      style={[
        styles.card,
        {
          backgroundColor: theme.colors.surface,
          borderColor: theme.colors.outline,
        },
        style,
      ]}
      {...props}
    />
  );
}

export const AppCard = Object.assign(AppCardBase, {
  Actions: Card.Actions,
  Content: Card.Content,
  Title: Card.Title,
});

const styles = StyleSheet.create({
  card: {
    borderWidth: 1,
    borderRadius: 18,
  },
});

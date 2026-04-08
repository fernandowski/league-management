import { PropsWithChildren } from 'react';
import { DimensionValue, StyleProp, StyleSheet, View, ViewStyle } from 'react-native';
import { Modal, Portal } from 'react-native-paper';

import { useAppTheme } from '@/theme/theme';

type AppModalProps = PropsWithChildren<{
  visible: boolean;
  onDismiss?: () => void;
  dismissable?: boolean;
  width?: DimensionValue;
  height?: DimensionValue;
  contentContainerStyle?: StyleProp<ViewStyle>;
}>;

export function AppModal({
  visible,
  onDismiss,
  dismissable = false,
  width,
  height,
  contentContainerStyle,
  children,
}: AppModalProps) {
  const theme = useAppTheme();

  return (
    <Portal>
      <Modal
        visible={visible}
        onDismiss={onDismiss}
        dismissable={dismissable}
        contentContainerStyle={[
          styles.modal,
          {
            backgroundColor: theme.colors.surface,
            borderColor: theme.colors.outline,
            shadowColor: theme.colors.shadow,
          },
          width ? { width } : styles.defaultWidth,
          height ? { height } : styles.defaultHeight,
          contentContainerStyle,
        ]}
      >
        <View style={styles.content}>{children}</View>
      </Modal>
    </Portal>
  );
}

const styles = StyleSheet.create({
  modal: {
    alignSelf: 'center',
    borderRadius: 24,
    borderWidth: 1,
    padding: 20,
    shadowOffset: { width: 0, height: 10 },
    shadowOpacity: 0.12,
    shadowRadius: 18,
    elevation: 10,
  },
  content: {
    gap: 16,
  },
  defaultWidth: {
    width: 'auto',
  },
  defaultHeight: {
    height: 'auto',
  },
});

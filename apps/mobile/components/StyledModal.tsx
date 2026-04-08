import { PropsWithChildren } from 'react';
import { DimensionValue, StyleProp, ViewStyle } from 'react-native';

import { AppModal } from '@/components/ui/AppModal';

interface ModalProps extends PropsWithChildren {
  isOpen: boolean;
  onDismiss?: () => void;
  dismissable?: boolean;
  width?: DimensionValue;
  height?: DimensionValue;
  contentContainerStyle?: StyleProp<ViewStyle>;
}

export default function StyledModal({
  isOpen,
  onDismiss,
  dismissable,
  width,
  height,
  contentContainerStyle,
  children,
}: ModalProps) {
  return (
    <AppModal
      visible={isOpen}
      onDismiss={onDismiss}
      dismissable={dismissable}
      width={width}
      height={height}
      contentContainerStyle={contentContainerStyle}
    >
      {children}
    </AppModal>
  );
}

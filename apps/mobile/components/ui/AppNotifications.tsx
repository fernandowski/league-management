import { createContext, PropsWithChildren, useContext, useEffect, useMemo, useState } from 'react';
import { Platform, StyleSheet, useWindowDimensions, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { useAppTheme } from '@/theme/theme';

type NotificationVariant = 'error' | 'warning' | 'info';

interface NotificationItem {
  id: string;
  message: string;
  variant: NotificationVariant;
}

interface NotificationContextValue {
  showNotification: (message: string, variant?: NotificationVariant) => void;
  dismissNotification: (id: string) => void;
}

const NotificationContext = createContext<NotificationContextValue | null>(null);

export function AppNotificationsProvider({ children }: PropsWithChildren) {
  const [notifications, setNotifications] = useState<NotificationItem[]>([]);
  const value = useMemo<NotificationContextValue>(() => ({
    showNotification: (message, variant = 'info') => {
      if (!message.trim()) {
        return;
      }

      const id = `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
      setNotifications((current) => [...current.slice(-2), { id, message, variant }]);
    },
    dismissNotification: (id) => {
      setNotifications((current) => current.filter((item) => item.id !== id));
    },
  }), []);

  return (
    <NotificationContext.Provider value={value}>
      {children}
      <NotificationViewport notifications={notifications} onDismiss={value.dismissNotification} />
    </NotificationContext.Provider>
  );
}

export function useNotification() {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error('useNotification must be used within AppNotificationsProvider');
  }

  return context;
}

function NotificationViewport({
  notifications,
  onDismiss,
}: {
  notifications: NotificationItem[];
  onDismiss: (id: string) => void;
}) {
  const dimensions = useWindowDimensions();
  const isCompact = dimensions.width < 768;

  return (
    <View pointerEvents="box-none" style={[styles.viewport, isCompact ? styles.viewportCompact : styles.viewportWide]}>
      {notifications.map((notification) => (
        <NotificationToast key={notification.id} notification={notification} onDismiss={onDismiss} compact={isCompact} />
      ))}
    </View>
  );
}

function NotificationToast({
  notification,
  onDismiss,
  compact,
}: {
  notification: NotificationItem;
  onDismiss: (id: string) => void;
  compact: boolean;
}) {
  const theme = useAppTheme();

  useEffect(() => {
    const timeout = setTimeout(() => onDismiss(notification.id), 4200);
    return () => clearTimeout(timeout);
  }, [notification.id, onDismiss]);

  const colors = useMemo(() => {
    switch (notification.variant) {
      case 'error':
        return {
          backgroundColor: theme.colors.errorContainer,
          borderColor: theme.colors.error,
          textColor: theme.colors.onErrorContainer,
        };
      case 'warning':
        return {
          backgroundColor: theme.colors.tertiaryContainer,
          borderColor: theme.colors.tertiary,
          textColor: theme.colors.onTertiaryContainer,
        };
      default:
        return {
          backgroundColor: theme.colors.primaryContainer,
          borderColor: theme.colors.primary,
          textColor: theme.colors.onPrimaryContainer,
        };
    }
  }, [notification.variant, theme.colors]);

  return (
    <View
      style={[
        styles.toast,
        compact ? styles.toastCompact : styles.toastWide,
        {
          backgroundColor: colors.backgroundColor,
          borderColor: colors.borderColor,
          shadowColor: theme.colors.shadow,
        },
      ]}
    >
      <AppText style={{ color: colors.textColor }}>{notification.message}</AppText>
    </View>
  );
}

const styles = StyleSheet.create({
  viewport: {
    position: 'absolute',
    top: Platform.select({ web: 24, default: 12 }),
    gap: 10,
    zIndex: 1000,
    elevation: 1000,
  },
  viewportCompact: {
    left: 12,
    right: 12,
    alignItems: 'stretch',
  },
  viewportWide: {
    right: 24,
    maxWidth: 420,
    width: '100%',
    alignItems: 'stretch',
  },
  toast: {
    borderWidth: 1,
    borderRadius: 16,
    paddingHorizontal: 14,
    paddingVertical: 12,
    shadowOpacity: 0.12,
    shadowRadius: 16,
    shadowOffset: { width: 0, height: 10 },
  },
  toastCompact: {
    width: '100%',
  },
  toastWide: {
    width: '100%',
  },
});

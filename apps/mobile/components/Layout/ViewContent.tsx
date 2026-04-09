import React from 'react';
import { ScrollView, StyleSheet, View } from 'react-native';

import { Header } from '@/components/Header/Header';
import { useAppTheme } from '@/theme/theme';

interface ViewContentProps {
  children: React.ReactNode;
}

const ViewContent = ({ children }: ViewContentProps) => {
  const theme = useAppTheme();

  return (
    <View style={[styles.outerContainer, { backgroundColor: theme.colors.background }]}>
      <View
        style={[
          styles.headerWrap,
          {
            backgroundColor: theme.colors.surface,
            borderBottomColor: theme.colors.outlineVariant,
            shadowColor: theme.colors.shadow,
          },
        ]}
      >
        <Header />
      </View>

      <ScrollView
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
      >
        {children}
      </ScrollView>
    </View>
  );
};

const styles = StyleSheet.create({
  outerContainer: {
    flex: 1,
  },
  headerWrap: {
    paddingHorizontal: 20,
    paddingBottom: 4,
    borderBottomWidth: 1,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.06,
    shadowRadius: 10,
    elevation: 2,
  },
  scrollView: {
    flex: 1,
    width: '100%',
  },
  scrollContent: {
    paddingHorizontal: 20,
    paddingTop: 24,
    paddingBottom: 24,
  },
});

export default ViewContent;

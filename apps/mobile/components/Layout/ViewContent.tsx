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
      <View style={styles.headerWrap}>
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
    paddingTop: 20,
  },
  scrollView: {
    flex: 1,
    width: '100%',
  },
  scrollContent: {
    paddingHorizontal: 20,
    paddingTop: 12,
    paddingBottom: 24,
  },
});

export default ViewContent;

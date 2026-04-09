import React, { useState } from 'react';
import { StyleSheet, useWindowDimensions, View } from 'react-native';

import { AppButton } from '@/components/ui/AppButton';
import { AppText } from '@/components/ui/AppText';
import { useAppTheme } from '@/theme/theme';

type Tab = {
  key: string;
  title: string;
  view: React.ReactNode;
};

interface TabsProps {
  tabs: Tab[];
}

const Tabs = ({ tabs }: TabsProps): React.ReactNode => {
  const dimensions = useWindowDimensions();
  const isLargeScreen = dimensions.width >= 768;
  const [activeTab, setActiveTab] = useState(0);
  const theme = useAppTheme();

  return (
    <View style={styles.container}>
      <View
        style={[
          styles.tabContainer,
          {
            backgroundColor: theme.colors.surface,
            borderColor: theme.colors.outline,
          },
          !isLargeScreen && styles.tabContainerStacked,
        ]}
      >
        {tabs.map(({ key, title }, index) => (
          <AppButton
            key={key}
            mode={activeTab === index ? 'contained' : 'contained-tonal'}
            style={[
              styles.tab,
              {
                backgroundColor: activeTab === index ? theme.colors.primary : theme.colors.surfaceVariant,
              },
            ]}
            onPress={() => setActiveTab(index)}
          >
            <AppText
              style={[
                styles.tabText,
                {
                  color:
                    activeTab === index ? theme.colors.onPrimary : theme.colors.onSurfaceVariant,
                },
              ]}
            >
              {title}
            </AppText>
          </AppButton>
        ))}
      </View>
      <View style={styles.content}>{tabs[activeTab].view}</View>
    </View>
  );
};

export default Tabs;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    gap: 16,
  },
  tabContainer: {
    flexDirection: 'row',
    borderWidth: 1,
    borderRadius: 14,
    padding: 4,
    gap: 4,
    alignSelf: 'flex-start',
    flexWrap: 'wrap',
  },
  tabContainerStacked: {
    flexDirection: 'column',
    alignSelf: 'stretch',
  },
  tab: {
    borderRadius: 10,
    alignItems: 'center',
  },
  tabText: {
    fontSize: 14,
    fontWeight: '600',
  },
  content: {
    flex: 1,
  },
});
